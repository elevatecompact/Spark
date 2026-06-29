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
	"github.com/elevatecompact/spark/services/payment-service/internal/config"
	"github.com/elevatecompact/spark/services/payment-service/internal/domain"
	"github.com/elevatecompact/spark/services/payment-service/internal/events"
	"github.com/elevatecompact/spark/services/payment-service/internal/handler"
	"github.com/elevatecompact/spark/services/payment-service/internal/processor"
	"github.com/elevatecompact/spark/services/payment-service/internal/repository"
	"github.com/elevatecompact/spark/services/payment-service/internal/service"
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

	intentRepo := repository.NewPaymentIntentRepository(pool)
	methodRepo := repository.NewPaymentMethodRepository(pool)
	webhookRepo := repository.NewWebhookRepository(pool)

	processors := map[domain.PaymentProcessor]processor.PaymentProcessor{
		domain.ProcessorStripe: processor.NewNoopProcessor(domain.ProcessorStripe),
		domain.ProcessorPayPal: processor.NewNoopProcessor(domain.ProcessorPayPal),
	}

	paymentSvc := service.NewPaymentService(
		intentRepo, methodRepo, webhookRepo, processors, eventProducer,
		cfg.AppConfig.StripeEnabled, cfg.AppConfig.PayPalEnabled,
		cfg.AppConfig.SaveMethods, cfg.AppConfig.RefundsEnabled,
	)

	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/v1/payment-intents", func(r chi.Router) {
		r.Post("/", paymentHandler.CreateIntent)
		r.Get("/", paymentHandler.ListIntents)
		r.Get("/{id}", paymentHandler.GetIntent)
		r.Post("/{id}/confirm", paymentHandler.ConfirmIntent)
		r.Post("/{id}/cancel", paymentHandler.CancelIntent)
		r.Post("/{id}/refund", paymentHandler.RefundIntent)
	})
	r.Route("/v1/payment-methods", func(r chi.Router) {
		r.Post("/", paymentHandler.CreatePaymentMethod)
		r.Get("/", paymentHandler.ListPaymentMethods)
		r.Get("/{id}", paymentHandler.GetPaymentMethod)
		r.Patch("/{id}", paymentHandler.SetDefaultPaymentMethod)
		r.Delete("/{id}", paymentHandler.DeletePaymentMethod)
	})
	r.Get("/v1/refunds", paymentHandler.ListRefunds)
	r.Post("/v1/payouts", paymentHandler.CreatePayout)
	r.Get("/v1/payouts/{id}", paymentHandler.GetPayout)
	r.Post("/v1/webhooks/stripe", paymentHandler.ProcessStripeWebhook)
	r.Post("/v1/webhooks/paypal", paymentHandler.ProcessPayPalWebhook)
	r.Get("/v1/admin/processors/status", paymentHandler.GetProcessorStatus)
	r.Post("/v1/admin/webhooks/retry/{id}", paymentHandler.RetryWebhook)

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
		log.Info().Str("port", cfg.Service.Port).Msg("starting payment-service")
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
