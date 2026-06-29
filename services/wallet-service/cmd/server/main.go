package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/elevatecompact/spark/packages/database"
	"github.com/elevatecompact/spark/services/wallet-service/config"
	"github.com/elevatecompact/spark/services/wallet-service/internal/events"
	"github.com/elevatecompact/spark/services/wallet-service/internal/handler"
	"github.com/elevatecompact/spark/services/wallet-service/internal/repository"
	"github.com/elevatecompact/spark/services/wallet-service/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("starting wallet-service")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	pgCfg := database.DefaultPGConfig(pgURL)

	bgCtx := context.Background()
	pool, err := database.NewPool(bgCtx, pgCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	defer database.Close(pool)
	log.Info().Msg("connected to postgres")

	if err := database.RunMigrations(bgCtx, pool, "services/wallet-service/migrations"); err != nil {
		log.Warn().Err(err).Msg("failed to run migrations")
	}

	redisURL := fmt.Sprintf("redis://%s", cfg.RedisAddr)
	if cfg.RedisPassword != "" {
		redisURL = fmt.Sprintf("redis://:%s@%s", cfg.RedisPassword, cfg.RedisAddr)
	}
	rCfg := database.DefaultRedisConfig(redisURL)
	rCfg.DB = cfg.RedisDB

	rdb, err := database.NewRedisClient(rCfg)
	if err != nil {
		log.Warn().Err(err).Msg("redis not available, creating dummy client")
		rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	} else {
		defer database.CloseRedis(rdb)
		log.Info().Msg("connected to redis")
	}

	var eventPub events.EventProducer
	if len(cfg.KafkaBrokers) > 0 && cfg.KafkaBrokers[0] != "" {
		eventPub = events.NewKafkaProducer(cfg.KafkaBrokers)
		log.Info().Strs("brokers", cfg.KafkaBrokers).Msg("kafka producer created")
	} else {
		eventPub = events.NewNoopProducer()
		log.Info().Msg("using noop event producer")
	}

	walletRepo := repository.NewWalletRepository(pool)
	txnRepo := repository.NewTransactionRepository(pool)
	payoutRepo := repository.NewPayoutRepository(pool)

	walletSvc := service.NewWalletService(walletRepo, eventPub, cfg.MaxBalanceCents)
	txnSvc := service.NewTransactionService(walletRepo, txnRepo, eventPub, cfg.MaxBalanceCents)
	payProc := service.NewNoopPaymentProcessor()
	payoutSvc := service.NewPayoutService(walletRepo, payoutRepo, eventPub, payProc, cfg.PayoutMinCents)

	walletH := handler.NewWalletHandler(walletSvc)
	txnH := handler.NewTransactionHandler(txnSvc)
	payoutH := handler.NewPayoutHandler(payoutSvc)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(corsMiddleware(cfg.AllowedOrigins))
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/wallets", func(r chi.Router) {
			r.Get("/me", walletH.GetMyWallet)
			r.Get("/me/balances", walletH.GetMyBalances)
			r.Get("/by-user", walletH.GetByUserID)
			r.Post("/freeze", walletH.Freeze)
			r.Post("/close", walletH.Close)
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Post("/deposit", txnH.Deposit)
			r.Post("/withdraw", txnH.Withdraw)
			r.Post("/transfer", txnH.Transfer)
			r.Post("/tip", txnH.Tip)
			r.Get("/", txnH.List)
			r.Get("/{id}", txnH.Get)
		})

		r.Route("/payouts", func(r chi.Router) {
			r.Post("/request", payoutH.Request)
			r.Get("/", payoutH.List)
			r.Get("/{id}", payoutH.Get)
		})
	})

	wsSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	ctx, cancel := context.WithCancel(context.Background())
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		log.Info().Int("port", cfg.Port).Msg("http server listening")
		if err := wsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("http server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-sigCh:
			log.Info().Str("signal", sig.String()).Msg("received signal")
			cancel()
		case <-gCtx.Done():
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		log.Info().Msg("shutting down http server")
		return wsSrv.Shutdown(shutdownCtx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}

	if err := eventPub.Close(); err != nil {
		log.Warn().Err(err).Msg("failed to close event producer")
	}

	log.Info().Msg("wallet-service stopped")
}

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	originMap := make(map[string]bool)
	for _, o := range allowedOrigins {
		originMap[o] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if originMap["*"] || originMap[origin] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
