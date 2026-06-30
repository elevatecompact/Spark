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
	"github.com/elevatecompact/spark/services/analytics-service/internal/config"
	"github.com/elevatecompact/spark/services/analytics-service/internal/events"
	"github.com/elevatecompact/spark/services/analytics-service/internal/handler"
	"github.com/elevatecompact/spark/services/analytics-service/internal/repository"
	"github.com/elevatecompact/spark/services/analytics-service/internal/service"
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

	// Open ClickHouse if a DSN is configured. A nil connection falls back to
	// Postgres for realtime metrics and queries.
	chConn, err := repository.OpenClickHouse(ctx, cfg.ClickHouse.URL)
	if err != nil {
		log.Warn().Err(err).Msg("clickhouse disabled; falling back to postgres")
		chConn = nil
	} else if chConn != nil {
		defer chConn.Close()
		log.Info().Msg("clickhouse connected")
	} else {
		log.Info().Msg("clickhouse not configured; using postgres for realtime metrics")
	}

	var eventProducer events.EventProducer
	var eventConsumer events.EventConsumer
	if cfg.Kafka.Enabled {
		eventProducer = events.NewKafkaProducer(cfg.Kafka.Brokers)
		eventConsumer = events.NewKafkaConsumer(cfg.Kafka.Brokers, cfg.Kafka.GroupID)
	} else {
		eventProducer = events.NewNoopProducer()
		eventConsumer = events.NewNoopConsumer()
	}
	defer eventProducer.Close()

	eventRepo := repository.NewTrackedEventRepository(pool, chConn)
	dashRepo := repository.NewDashboardRepository(pool)
	reportRepo := repository.NewReportRepository(pool)
	tmplRepo := repository.NewReportTemplateRepository(pool)
	funnelRepo := repository.NewFunnelRepository(pool)

	analyticsSvc := service.NewAnalyticsService(eventRepo, dashRepo, reportRepo, tmplRepo, funnelRepo, eventProducer)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsSvc)

	go func() {
		if err := eventConsumer.Consume(ctx, func(ctx context.Context, event events.CloudEvent) error {
			log.Debug().Str("type", event.Type).Msg("consumed event")
			return nil
		}); err != nil {
			log.Warn().Err(err).Msg("event consumer stopped")
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Post("/v1/events/track", analyticsHandler.TrackEvent)
	r.Post("/v1/events/batch", analyticsHandler.TrackBatch)

	r.Route("/v1/dashboards", func(r chi.Router) {
		r.Get("/creator/{id}", analyticsHandler.GetCreatorDashboard)
		r.Get("/viewer/{id}", analyticsHandler.GetViewerDashboard)
		r.Get("/admin", analyticsHandler.GetAdminDashboard)
	})

	r.Route("/v1/metrics", func(r chi.Router) {
		r.Get("/realtime", analyticsHandler.GetRealtimeMetrics)
		r.Get("/historical", analyticsHandler.GetHistoricalMetrics)
		r.Post("/query", analyticsHandler.QueryMetrics)
	})

	r.Route("/v1/reports", func(r chi.Router) {
		r.Post("/generate", analyticsHandler.GenerateReport)
		r.Get("/", analyticsHandler.ListReports)
		r.Get("/{id}", analyticsHandler.GetReport)
		r.Get("/templates", analyticsHandler.ListTemplates)
	})

	r.Route("/v1/funnels", func(r chi.Router) {
		r.Post("/define", analyticsHandler.DefineFunnel)
		r.Get("/{id}", analyticsHandler.GetFunnel)
		r.Get("/{id}/analyze", analyticsHandler.AnalyzeFunnel)
		r.Get("/{id}/results", analyticsHandler.AnalyzeFunnel)
	})

	r.Post("/v1/export/csv", analyticsHandler.ExportCSV)
	r.Post("/v1/export/json", analyticsHandler.ExportJSON)

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
		log.Info().Str("port", cfg.Service.Port).Msg("starting analytics-service")
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
