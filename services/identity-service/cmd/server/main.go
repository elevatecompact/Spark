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

	"github.com/elevatecompact/spark/services/identity-service/internal/config"
	"github.com/elevatecompact/spark/services/identity-service/internal/database"
	"github.com/elevatecompact/spark/services/identity-service/internal/events"
	"github.com/elevatecompact/spark/services/identity-service/internal/handler"
	"github.com/elevatecompact/spark/services/identity-service/internal/repository"
	"github.com/elevatecompact/spark/services/identity-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	setupLogging(cfg.LogLevel)

	log.Info().Str("env", cfg.Env).Msg("starting identity service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	if err := db.RunMigrations(ctx, "migrations"); err != nil {
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
		eventPub = events.NewKafkaProducer(cfg.KafkaBrokers, "identity-events")
		log.Info().Str("brokers", cfg.KafkaBrokers[0]).Msg("initialized kafka producer")
	} else {
		eventPub = events.NewNoopProducer()
		log.Info().Msg("using noop event producer")
	}
	defer eventPub.Close()

	userRepo := repository.NewUserRepository(db.Pool)
	sessionRepo := repository.NewSessionRepository(redisClient, cfg.JWTExpiry)
	oauthRepo := repository.NewOAuthRepository(db.Pool)

	tokenSvc := service.NewTokenService(cfg.JWTSecret, cfg.JWTExpiry)
	authSvc := service.NewAuthService(userRepo, sessionRepo, tokenSvc, eventPub, 24*time.Hour)
	userSvc := service.NewUserService(userRepo, eventPub)
	oauthCfg := service.OAuthConfig{
		AccessTokenTTL:  cfg.JWTExpiry,
		RefreshTokenTTL: 7 * 24 * time.Hour,
		AuthCodeTTL:     10 * time.Minute,
	}
	oauthSvc := service.NewOAuthService(oauthRepo, userRepo, tokenSvc, authSvc, oauthCfg)
	passkeyCfg := service.PasskeyConfig{
		RPID:   cfg.PasskeyRPID,
		RPName: cfg.PasskeyRPName,
		Origin: cfg.PasskeyOrigin,
	}
	passkeySvc := service.NewPasskeyService(sessionRepo, tokenSvc, passkeyCfg)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	oauthHandler := handler.NewOAuthHandler(oauthSvc)
	passkeyHandler := handler.NewPasskeyHandler(passkeySvc)

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

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)

		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(tokenSvc, authSvc))
			r.Post("/logout", authHandler.Logout)
			r.Post("/logout-all", authHandler.LogoutAll)
		})
	})

	r.Route("/oauth", func(r chi.Router) {
		r.Get("/authorize", oauthHandler.AuthorizeForm)
		r.Post("/token", oauthHandler.Token)
		r.Post("/introspect", oauthHandler.Introspect)
		r.Get("/userinfo", oauthHandler.UserInfo)
		r.Get("/.well-known/openid-configuration", oauthHandler.OpenIDDiscovery)

		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(tokenSvc, authSvc))
			r.Post("/authorize", oauthHandler.Authorize)
		})
	})

	r.Route("/passkeys", func(r chi.Router) {
		r.Post("/authenticate", passkeyHandler.BeginAuthentication)
		r.Post("/authenticate/complete", passkeyHandler.FinishAuthentication)

		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(tokenSvc, authSvc))
			r.Post("/register", passkeyHandler.BeginRegistration)
			r.Post("/register/complete", passkeyHandler.FinishRegistration)
			r.Get("/", passkeyHandler.ListPasskeys)
			r.Delete("/{id}", passkeyHandler.DeletePasskey)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(tokenSvc, authSvc))
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateMe)
			r.Post("/me/change-password", userHandler.ChangePassword)
			r.Delete("/me", userHandler.DeleteAccount)
		})

		r.Group(func(r chi.Router) {
			r.Use(handler.AuthRequired(tokenSvc, authSvc))
			r.Use(handler.ModOrAdminRequired)
			r.Get("/", userHandler.SearchUsers)
			r.Get("/{id}", userHandler.GetUserByID)
			r.Post("/{id}/verify", userHandler.VerifyUser)
			r.Post("/{id}/suspend", userHandler.SuspendUser)
			r.Put("/{id}/role", userHandler.UpdateRole)
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		healthErr := db.HealthCheck(r.Context())
		redisErr := redisClient.Ping(r.Context()).Err()

		status := http.StatusOK
		if healthErr != nil || redisErr != nil {
			status = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		healthStatus := map[string]string{
			"status":  "ok",
			"service": "identity-service",
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
