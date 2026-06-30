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
	"github.com/elevatecompact/spark/services/notification-service/internal/config"
	"github.com/elevatecompact/spark/services/notification-service/internal/events"
	"github.com/elevatecompact/spark/services/notification-service/internal/handler"
	"github.com/elevatecompact/spark/services/notification-service/internal/processor"
	"github.com/elevatecompact/spark/services/notification-service/internal/repository"
	"github.com/elevatecompact/spark/services/notification-service/internal/service"
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

	if err := database.RunMigrations(ctx, pool, "migrations"); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	var eventPub events.EventProducer
	if cfg.Kafka.Enabled {
		eventPub = events.NewKafkaProducer(cfg.Kafka.Brokers)
	} else {
		eventPub = events.NewNoopProducer()
	}
	defer eventPub.Close()

	notifRepo := repository.NewNotificationRepository(pool)
	prefRepo := repository.NewPreferenceRepository(pool)
	devRepo := repository.NewDeviceRepository(pool)
	tmplRepo := repository.NewTemplateRepository(pool)

	// Wire real processors. They automatically fall back to noop when
	// credentials are missing, so the service is still useful in dev.
	push := processor.NewFCMPush(cfg.Push.FCMKey)
	email := processor.NewSendGridEmail(cfg.Email.SendGridAPIKey)
	sms := processor.NewTwilioSMS(cfg.SMS.TwilioSID, cfg.SMS.TwilioToken, cfg.SMS.TwilioPhone)

	svc := service.NewNotificationService(
		notifRepo, prefRepo, devRepo, tmplRepo, push, email, sms, eventPub,
		cfg.AppConfig.PushEnabled, cfg.AppConfig.EmailEnabled, cfg.AppConfig.SMSEnabled, cfg.AppConfig.InAppEnabled,
	)

	h := handler.NewNotifHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/v1/notifications", func(r chi.Router) {
		r.Get("/", h.ListNotifications)
		r.Post("/read-all", h.MarkAllRead)
		r.Patch("/{id}/read", h.MarkRead)
		r.Delete("/{id}", h.Delete)
	})
	r.Get("/v1/preferences", h.GetPreferences)
	r.Patch("/v1/preferences", h.UpdatePreferences)
	r.Get("/v1/templates", h.ListTemplates)
	r.Post("/v1/templates", h.CreateTemplate)
	r.Patch("/v1/templates/{id}", h.UpdateTemplate)
	r.Post("/v1/send", h.SendNotification)
	r.Post("/v1/send/batch", h.SendBatch)
	r.Post("/v1/devices", h.RegisterDevice)
	r.Delete("/v1/devices/{id}", h.UnregisterDevice)
	r.Get("/v1/devices", h.ListDevices)
	r.Post("/v1/admin/test-push", h.TestPush)
	r.Post("/v1/admin/test-email", h.TestEmail)
	r.Get("/v1/admin/delivery-stats", h.DeliveryStats)

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
		log.Info().Str("port", cfg.Service.Port).Msg("starting notification-service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	<-quit
	log.Info().Msg("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
	log.Info().Msg("server exited")
}
