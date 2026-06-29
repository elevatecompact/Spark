package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/packages/database"
	"github.com/elevatecompact/spark/services/translation-service/internal/config"
	"github.com/elevatecompact/spark/services/translation-service/internal/events"
	"github.com/elevatecompact/spark/services/translation-service/internal/handler"
	"github.com/elevatecompact/spark/services/translation-service/internal/repository"
	"github.com/elevatecompact/spark/services/translation-service/internal/service"
)

func main() {
	cfg := config.Load()
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

	repo := repository.NewTranslationRepository(pool)

	var eventPub events.EventProducer
	if len(cfg.KafkaBrokers) > 0 {
		eventPub = events.NewKafkaProducer(cfg.KafkaBrokers)
	} else {
		eventPub = events.NewNoopProducer()
	}

	svc := service.NewTranslationService(repo, eventPub)
	h := handler.New(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"*"}}))

	h.Register(r)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Port).Msg("translation-service starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("forced shutdown")
	}
	eventPub.Close()
	log.Info().Msg("server stopped")
}
