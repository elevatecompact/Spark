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
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/elevatecompact/spark/services/stream-service/internal/config"
	"github.com/elevatecompact/spark/services/stream-service/internal/database"
	"github.com/elevatecompact/spark/services/stream-service/internal/events"
	"github.com/elevatecompact/spark/services/stream-service/internal/handler"
	"github.com/elevatecompact/spark/services/stream-service/internal/repository"
	"github.com/elevatecompact/spark/services/stream-service/internal/service"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Logger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	log.Info().Int("port", cfg.Server.Port).Msg("Starting stream service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.NewPostgresDB(
		ctx,
		cfg.DatabaseDSN(),
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to PostgreSQL")
	}
	defer db.Close()

	if err := db.RunMigrations(ctx, "migrations"); err != nil {
		log.Warn().Err(err).Msg("Failed to run database migrations")
	}

	producer := events.NewEventProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	defer producer.Close()

	streamRepo := repository.NewStreamRepository(db.Pool)

	hub := service.NewWebSocketHub()

	streamSvc := service.NewStreamService(streamRepo, producer, cfg, hub)
	rtmpSvc := service.NewRTMPService(streamRepo, streamSvc)
	webrtcSvc := service.NewWebRTCService(cfg, streamSvc)
	hlsSvc := service.NewHLSService(cfg, streamSvc)
	healthSvc := service.NewHealthService()
	recordingSvc := service.NewRecordingService(cfg, producer, streamRepo)

	mw := handler.NewMiddleware(cfg)

	streamH := handler.NewStreamHandler(streamSvc, rtmpSvc)
	webrtcH := handler.NewWebRTCHandler(webrtcSvc, streamSvc)
	playbackH := handler.NewPlaybackHandler(hlsSvc, streamSvc)

	r := chi.NewRouter()

	r.Use(mw.Recovery)
	r.Use(mw.Logger)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/streams", func(r chi.Router) {
			r.Post("/", streamH.CreateStream)
			r.Get("/", streamH.ListStreams)
			r.Get("/live", streamH.ListLiveStreams)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", streamH.GetStream)
				r.Put("/", streamH.UpdateStream)
				r.Post("/end", streamH.EndStream)
				r.Get("/health", streamH.GetStreamHealth)

				r.Get("/watch", func(w http.ResponseWriter, r *http.Request) {
					handleWebSocket(w, r, hub, streamSvc)
				})
			})
		})

		r.Route("/webrtc", func(r chi.Router) {
			r.Route("/stream/{id}", func(r chi.Router) {
				r.Post("/offer", webrtcH.HandleOffer)
				r.Post("/answer", webrtcH.HandleAnswer)
				r.Post("/ice-candidate", webrtcH.HandleICECandidate)
				r.Post("/viewer/join", webrtcH.JoinViewer)
				r.Post("/viewer/leave", webrtcH.LeaveViewer)
				r.Get("/", webrtcH.GetStream)
			})
		})

		r.Route("/playback/{id}", func(r chi.Router) {
			r.Get("/info", playbackH.GetStreamInfo)
			r.Get("/hls/master.m3u8", playbackH.GetMasterPlaylist)
			r.Get("/hls/{quality}/index.m3u8", playbackH.GetQualityPlaylist)
			r.Get("/hls/{quality}/{segment}", playbackH.GetSegment)
			r.Get("/thumbnails/{timestamp}.jpg", playbackH.GetThumbnail)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		count, _ := streamSvc.GetLiveCount(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","service":"stream-service","live_streams":%d}`, count)
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("addr", srv.Addr).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	<-quit
	log.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.GracefulTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited gracefully")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, hub *service.WebSocketHub, streamSvc *service.StreamService) {
	streamIDStr := chi.URLParam(r, "id")
	streamID, err := uuid.Parse(streamIDStr)
	if err != nil {
		http.Error(w, "invalid stream id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Str("stream_id", streamID.String()).Msg("WebSocket upgrade failed")
		return
	}

	client := &service.WebSocketClient{
		ID:       uuid.New(),
		StreamID: streamID,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	hub.Register(client)

	go clientWritePump(client)
	go clientReadPump(client, hub, streamSvc)
}

func clientWritePump(client *service.WebSocketClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func clientReadPump(client *service.WebSocketClient, hub *service.WebSocketHub, streamSvc *service.StreamService) {
	defer func() {
		hub.Unregister(client)
		client.Conn.Close()
	}()

	client.Conn.SetReadLimit(4096)
	client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.Conn.SetPongHandler(func(string) error {
		client.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := client.Conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
