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
	"github.com/elevatecompact/spark/services/chat-service/config"
	"github.com/elevatecompact/spark/services/chat-service/internal/events"
	"github.com/elevatecompact/spark/services/chat-service/internal/handler"
	"github.com/elevatecompact/spark/services/chat-service/internal/repository"
	"github.com/elevatecompact/spark/services/chat-service/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("starting chat-service")

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

	if err := database.RunMigrations(bgCtx, pool, "services/chat-service/migrations"); err != nil {
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
	if rdb != nil && len(cfg.KafkaBrokers) > 0 && cfg.KafkaBrokers[0] != "" {
		eventPub = events.NewKafkaProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
		log.Info().Str("topic", cfg.KafkaTopic).Strs("brokers", cfg.KafkaBrokers).Msg("kafka producer created")
	} else {
		eventPub = events.NewNoopProducer()
		log.Info().Msg("using noop event producer")
	}

	roomRepo := repository.NewRoomRepository(pool)
	msgRepo := repository.NewMessageRepository(pool)
	modRepo := repository.NewModerationRepository(pool, rdb)
	emoteRepo := repository.NewEmoteRepository(pool)

	hub := service.NewWebSocketHub()

	roomSvc := service.NewRoomService(roomRepo, eventPub)
	msgSvc := service.NewMessageService(msgRepo, roomRepo, modRepo, eventPub, cfg.MaxMessageLength)
	modSvc := service.NewModerationService(modRepo, eventPub)
	emoteSvc := service.NewEmoteService(emoteRepo)

	roomH := handler.NewRoomHandler(roomSvc)
	msgH := handler.NewMessageHandler(msgSvc)
	modH := handler.NewModerationHandler(modSvc)
	emoteH := handler.NewEmotesHandler(emoteSvc)
	wsH := handler.NewWebSocketHandler(hub)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(corsMiddleware(cfg.AllowedOrigins))
	r.Use(chiMiddleware.Timeout(30 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", roomH.Create)
			r.Get("/{id}", roomH.Get)
			r.Delete("/{id}", roomH.Close)

			r.Route("/{roomId}/messages", func(r chi.Router) {
				r.Get("/", msgH.GetHistory)
				r.Post("/", msgH.Send)
			})

			r.Route("/{roomId}/emotes", func(r chi.Router) {
				r.Get("/", emoteH.GetByRoom)
			})

			r.Route("/{roomId}/moderation", func(r chi.Router) {
				r.Post("/mute", modH.MuteUser)
				r.Post("/unmute/{userId}", modH.UnmuteUser)
				r.Post("/ban", modH.BanUser)
				r.Post("/unban/{userId}", modH.UnbanUser)
				r.Post("/slow-mode", modH.SetSlowMode)
			})
		})

		r.Put("/messages/{id}", msgH.Edit)
		r.Delete("/messages/{id}", msgH.Delete)

		r.Get("/emotes/global", emoteH.GetGlobal)
	})

	r.Get("/ws/chat/{roomId}", wsH.HandleWS)

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

	log.Info().Msg("chat-service stopped")
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
