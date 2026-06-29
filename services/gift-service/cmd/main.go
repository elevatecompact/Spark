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
	"github.com/elevatecompact/spark/services/gift-service/internal/config"
	"github.com/elevatecompact/spark/services/gift-service/internal/events"
	"github.com/elevatecompact/spark/services/gift-service/internal/handler"
	"github.com/elevatecompact/spark/services/gift-service/internal/repository"
	"github.com/elevatecompact/spark/services/gift-service/internal/service"
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

	itemRepo := repository.NewGiftItemRepository(pool)
	giftRepo := repository.NewGiftRepository(pool)
	cardRepo := repository.NewGiftCardRepository(pool)
	campaignRepo := repository.NewGiftCampaignRepository(pool)

	giftSvc := service.NewGiftService(giftRepo, itemRepo, cardRepo, campaignRepo, eventProducer, cfg.ToGiftServiceConfig())

	giftHandler := handler.NewGiftHandler(giftSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Route("/v1/gift-items", func(r chi.Router) {
		r.Post("/", giftHandler.CreateGiftItem)
		r.Get("/", giftHandler.ListGiftItems)
		r.Get("/{id}", giftHandler.GetGiftItem)
		r.Patch("/{id}", giftHandler.UpdateGiftItem)
		r.Delete("/{id}", giftHandler.DeleteGiftItem)
	})

	r.Route("/v1/gifts", func(r chi.Router) {
		r.Post("/send", giftHandler.SendGift)
		r.Post("/batch", giftHandler.SendBatchGift)
		r.Post("/subscription", giftHandler.SendSubscriptionGift)
		r.Get("/sent", giftHandler.ListSent)
		r.Get("/received", giftHandler.ListReceived)
		r.Post("/{id}/refund", giftHandler.RefundGift)
		r.Get("/{id}", giftHandler.GetGift)
		r.Get("/{id}/status", giftHandler.GetGift)
	})

	r.Post("/v1/gift-cards/purchase", giftHandler.PurchaseGiftCard)
	r.Post("/v1/gift-cards/redeem", giftHandler.RedeemGiftCard)
	r.Get("/v1/gift-cards/{code}", giftHandler.GetGiftCardByCode)

	r.Post("/v1/campaigns", giftHandler.CreateCampaign)
	r.Get("/v1/campaigns", giftHandler.ListCampaigns)
	r.Post("/v1/campaigns/{id}/match", giftHandler.ApplyCampaignMatch)

	r.Get("/v1/analytics/top-gifts", giftHandler.GetTopGifts)
	r.Get("/v1/analytics/leaderboard", giftHandler.GetLeaderboard)

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
		log.Info().Str("port", cfg.Service.Port).Msg("starting gift-service")
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
