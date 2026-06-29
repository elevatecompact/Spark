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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/packages/database"
	"github.com/elevatecompact/spark/services/subscription-service/internal/config"
	"github.com/elevatecompact/spark/services/subscription-service/internal/events"
	"github.com/elevatecompact/spark/services/subscription-service/internal/handler"
	"github.com/elevatecompact/spark/services/subscription-service/internal/repository"
	"github.com/elevatecompact/spark/services/subscription-service/internal/service"
)

func main() {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx := context.Background()

	dbCfg := database.DefaultPGConfig(cfg.Database.URL)
	pool, err := database.NewPool(ctx, dbCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer database.Close(pool)

	migrationsDir := "migrations"
	if err := database.RunMigrations(ctx, pool, migrationsDir); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	var eventProducer events.EventProducer
	if cfg.Kafka.Enabled {
		eventProducer = events.NewKafkaProducer(cfg.Kafka.Brokers)
	} else {
		eventProducer = events.NewNoopProducer()
	}
	defer eventProducer.Close()

	planRepo := repository.NewPlanRepository(pool)
	subRepo := repository.NewSubscriptionRepository(pool)
	invRepo := repository.NewInvoiceRepository(pool)

	planSvc := service.NewPlanService(planRepo)
	subSvc := service.NewSubscriptionService(subRepo, planRepo, invRepo, eventProducer,
		cfg.AppConfig.GraceDays, cfg.AppConfig.MaxActive, cfg.AppConfig.TrialDays)

	planHandler := handler.NewPlanHandler(planSvc)
	subHandler := handler.NewSubscriptionHandler(subSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/api/v1/plans", func(r chi.Router) {
		r.Post("/", planHandler.Create)
		r.Get("/", planHandler.List)
		r.Get("/{planId}", planHandler.GetByID)
		r.Put("/{planId}", planHandler.Update)
		r.Delete("/{planId}", planHandler.Delete)
	})

	r.Route("/api/v1/subscriptions", func(r chi.Router) {
		r.Post("/", subHandler.Subscribe)
		r.Get("/mine", subHandler.GetMy)
		r.Get("/{subId}", subHandler.GetByID)
		r.Post("/{subId}/cancel", subHandler.Cancel)
		r.Post("/{subId}/reactivate", subHandler.Reactivate)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Service.Port),
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.Service.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Service.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Service.IdleTimeout) * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("port", cfg.Service.Port).Msg("starting subscription-service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-quit
	log.Info().Msg("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited")
}
