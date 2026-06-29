package main

import (
	"context"
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

	"github.com/elevatecompact/spark/services/creator-service/internal/config"
	"github.com/elevatecompact/spark/services/creator-service/internal/database"
	"github.com/elevatecompact/spark/services/creator-service/internal/events"
	"github.com/elevatecompact/spark/services/creator-service/internal/handler"
	"github.com/elevatecompact/spark/services/creator-service/internal/repository"
	"github.com/elevatecompact/spark/services/creator-service/internal/service"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPostgresPool(ctx, cfg.Database.URL, database.PoolConfig{
		MaxOpenConnections:    cfg.Database.MaxOpenConnections,
		MaxIdleConnections:    cfg.Database.MaxIdleConnections,
		ConnectionMaxLifetime: cfg.Database.ConnectionMaxLifetime,
		ConnectionMaxIdleTime: cfg.Database.ConnectionMaxIdleTime,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer pool.Close()

	if err := database.RunMigrations(ctx, pool, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("Failed to run migrations")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Warn().Err(err).Msg("Redis not available, continuing without cache")
	}

	producer := events.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer producer.Close()

	creatorRepo := repository.NewCreatorRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	portfolioRepo := repository.NewPortfolioRepository(pool)

	cacheService := &redisCache{client: rdb}

	creatorSvc := service.NewCreatorService(creatorRepo, producer, cacheService)
	discoverSvc := service.NewDiscoverService(creatorRepo, categoryRepo)
	verifySvc := service.NewVerificationService(creatorRepo, producer)
	scheduleSvc := service.NewScheduleService(pool, creatorRepo)

	creatorH := handler.NewCreatorHandler(creatorSvc, verifySvc)
	categoryH := handler.NewCategoryHandler(categoryRepo, creatorSvc)
	discoverH := handler.NewDiscoverHandler(discoverSvc)
	scheduleH := handler.NewScheduleHandler(scheduleSvc)
	portfolioH := handler.NewPortfolioHandler(portfolioRepo)

	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(handler.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/creators", func(r chi.Router) {
			r.With(handler.OptionalAuth(cfg.JWT.Secret), handler.Pagination).Get("/", creatorH.Search)
			r.With(handler.AuthRequired(cfg.JWT.Secret)).Post("/", creatorH.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.With(handler.OptionalAuth(cfg.JWT.Secret)).Get("/", creatorH.GetByID)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Put("/", creatorH.Update)

				r.With(handler.AuthRequired(cfg.JWT.Secret)).Post("/follow", creatorH.Follow)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Post("/unfollow", creatorH.Unfollow)

				r.With(handler.Pagination).Get("/followers", creatorH.GetFollowers)
				r.With(handler.Pagination).Get("/following", creatorH.GetFollowing)

				r.With(handler.AdminRequired(cfg.JWT.Secret)).Post("/verify", creatorH.Verify)

				r.Get("/schedule", scheduleH.GetSchedule)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Post("/schedule", scheduleH.AddSlot)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Put("/schedule/{slotID}", scheduleH.UpdateSlot)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Delete("/schedule/{slotID}", scheduleH.DeleteSlot)

				r.Get("/portfolio", portfolioH.List)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Post("/portfolio", portfolioH.Create)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Put("/portfolio/{itemID}", portfolioH.Update)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Delete("/portfolio/{itemID}", portfolioH.Delete)
				r.With(handler.AuthRequired(cfg.JWT.Secret)).Put("/portfolio/{itemID}/featured", portfolioH.SetFeatured)
			})
		})

		r.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryH.List)
			r.With(handler.AdminRequired(cfg.JWT.Secret)).Post("/", categoryH.Create)
			r.Get("/{id}", categoryH.GetByID)
			r.With(handler.Pagination).Get("/{id}/creators", categoryH.GetCreators)
		})

		r.Route("/discover", func(r chi.Router) {
			r.With(handler.Pagination).Get("/trending", discoverH.Trending)
			r.With(handler.OptionalAuth(cfg.JWT.Secret), handler.Pagination).Get("/recommended", discoverH.Recommended)
			r.With(handler.Pagination).Get("/nearby", discoverH.Nearby)
			r.With(handler.Pagination).Get("/category/{categoryID}", discoverH.ByCategory)
			r.With(handler.Pagination).Get("/search", discoverH.Search)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"creator-service"}`))
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("Creator service starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}

type redisCache struct {
	client *redis.Client
}

func (c *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
