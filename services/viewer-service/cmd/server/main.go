package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/elevatecompact/spark/services/viewer-service/internal/config"
	"github.com/elevatecompact/spark/services/viewer-service/internal/handler"
	"github.com/elevatecompact/spark/services/viewer-service/internal/repository"
	"github.com/elevatecompact/spark/services/viewer-service/internal/service"
	"github.com/elevatecompact/spark/services/viewer-service/internal/events"
	"github.com/elevatecompact/spark/packages/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	setupLogging(cfg.LogLevel)

	log.Info().Str("env", cfg.Env).Msg("starting viewer service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPool(ctx, database.DefaultPGConfig(cfg.DatabaseURL))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer database.Close(pool)

	if err := database.RunMigrations(ctx, pool, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisURL,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("failed to connect to redis")
	}
	defer redisClient.Close()
	log.Info().Msg("connected to redis")

	var eventPub events.EventProducer
	if len(cfg.KafkaBrokers) > 0 && cfg.KafkaBrokers[0] != "" {
		eventPub = events.NewKafkaProducer(cfg.KafkaBrokers, "viewer-events")
		log.Info().Str("brokers", cfg.KafkaBrokers[0]).Msg("initialized kafka producer")
	} else {
		eventPub = events.NewNoopProducer()
		log.Info().Msg("using noop event producer")
	}
	defer eventPub.Close()

	watchHistoryRepo := repository.NewWatchHistoryRepository(pool)
	prefsRepo := repository.NewPreferencesRepository(pool)
	bookmarkRepo := repository.NewBookmarkRepository(pool)
	watchLaterRepo := repository.NewWatchLaterRepository(pool)
	ratingRepo := repository.NewRatingRepository(pool)
	reactionRepo := repository.NewReactionRepository(pool)
	reportRepo := repository.NewReportRepository(pool)

	watchSvc := service.NewWatchHistoryService(watchHistoryRepo, eventPub)
	prefsSvc := service.NewPreferencesService(prefsRepo)
	bookmarkSvc := service.NewBookmarkService(bookmarkRepo, cfg.MaxBookmarks)
	watchLaterSvc := service.NewWatchLaterService(watchLaterRepo, cfg.MaxWatchLater)
	engSvc := service.NewEngagementService(ratingRepo, reactionRepo, reportRepo, eventPub)

	watchHandler := handler.NewWatchHistoryHandler(watchSvc)
	prefsHandler := handler.NewPreferencesHandler(prefsSvc)
	bookmarkHandler := handler.NewBookmarkHandler(bookmarkSvc, watchLaterSvc)
	engHandler := handler.NewEngagementHandler(engSvc)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(handler.LoggerMiddleware)
	r.Use(handler.RequestIDMiddleware)
	r.Use(handler.RecoveryMiddleware)
	r.Use(chimiddleware.Timeout(30 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(handler.RateLimitMiddleware(100, 200))

	r.Route("/v1/history", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Post("/", watchHandler.RecordWatch)
			r.Get("/", watchHandler.GetHistory)
			r.Delete("/", watchHandler.ClearHistory)
			r.Delete("/{id}", watchHandler.DeleteEntry)
		})
	})

	r.Route("/v1/preferences", func(r chi.Router) {
		r.Get("/defaults", prefsHandler.GetDefault)
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Get("/", prefsHandler.Get)
			r.Put("/", prefsHandler.Replace)
			r.Patch("/", prefsHandler.Patch)
		})
	})

	r.Route("/v1/bookmarks", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Post("/", bookmarkHandler.Create)
			r.Get("/", bookmarkHandler.List)
			r.Delete("/{id}", bookmarkHandler.Delete)
		})
	})

	r.Route("/v1/watch-later", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Post("/", bookmarkHandler.AddWatchLater)
			r.Get("/", bookmarkHandler.ListWatchLater)
			r.Delete("/{id}", bookmarkHandler.RemoveWatchLater)
		})
	})

	r.Route("/v1/ratings", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Post("/", engHandler.RateContent)
		})
	})

	r.Route("/v1/reactions", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Post("/", engHandler.ToggleReaction)
		})
	})

	r.Route("/v1/reports", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(cfg.JWTSecret))
			r.Post("/", engHandler.ReportContent)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		healthErr := database.HealthCheck(r.Context(), pool)
		redisErr := redisClient.Ping(r.Context()).Err()

		status := http.StatusOK
		if healthErr != nil || redisErr != nil {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		healthStatus := map[string]string{
			"status":  "ok",
			"service": "viewer-service",
		}
		if healthErr != nil {
			healthStatus["database"] = "unhealthy"
		} else {
			healthStatus["database"] = "healthy"
		}
		if redisErr != nil {
			healthStatus["redis"] = "unhealthy"
		} else {
			healthStatus["redis"] = "healthy"
		}

		json.NewEncoder(w).Encode(healthStatus)
	})

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingUnaryInterceptor),
	)
	reflection.Register(grpcServer)

	go func() {
		log.Info().Str("port", cfg.ServerPort).Msg("http server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("http server failed")
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
		if err != nil {
			log.Fatal().Err(err).Msg("grpc server failed to listen")
		}
		log.Info().Str("port", cfg.GrpcPort).Msg("grpc server starting")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("grpc server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info().Str("signal", sig.String()).Msg("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	grpcServer.GracefulStop()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("http server shutdown error")
	}

	log.Info().Msg("server stopped gracefully")
}

func setupLogging(level string) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Caller().Logger()
}

func loggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	duration := time.Since(start)

	log.Debug().
		Str("method", info.FullMethod).
		Dur("duration", duration).
		Err(err).
		Msg("grpc request")

	return resp, err
}
