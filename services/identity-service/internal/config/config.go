package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	GrpcPort   string

	DatabaseURL string
	RedisURL    string

	JWTSecret     string
	JWTExpiry     time.Duration

	KafkaBrokers []string

	OAuthClientID     string
	OAuthClientSecret string

	AllowedOrigins []string

	PasskeyRPID     string
	PasskeyRPName   string
	PasskeyOrigin   string

	EncryptionKey string

	LogLevel string
	Env      string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("SERVER_PORT", "8080")
	v.SetDefault("GRPC_PORT", "9090")
	v.SetDefault("JWT_EXPIRY", "15m")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("ENV", "development")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	jwtExpiry, err := time.ParseDuration(v.GetString("JWT_EXPIRY"))
	if err != nil {
		jwtExpiry = 15 * time.Minute
	}

	cfg := &Config{
		ServerPort: v.GetString("SERVER_PORT"),
		GrpcPort:   v.GetString("GRPC_PORT"),

		DatabaseURL: v.GetString("DATABASE_URL"),
		RedisURL:    v.GetString("REDIS_URL"),

		JWTSecret: v.GetString("JWT_SECRET"),
		JWTExpiry: jwtExpiry,

		KafkaBrokers: strings.Split(v.GetString("KAFKA_BROKERS"), ","),

		OAuthClientID:     v.GetString("OAUTH_CLIENT_ID"),
		OAuthClientSecret: v.GetString("OAUTH_CLIENT_SECRET"),

		AllowedOrigins: strings.Split(v.GetString("ALLOWED_ORIGINS"), ","),

		PasskeyRPID:   v.GetString("PASSKEY_RP_ID"),
		PasskeyRPName: v.GetString("PASSKEY_RP_NAME"),
		PasskeyOrigin: v.GetString("PASSKEY_ORIGIN"),

		EncryptionKey: v.GetString("ENCRYPTION_KEY"),

		LogLevel: v.GetString("LOG_LEVEL"),
		Env:      v.GetString("ENV"),
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "change-me-in-production"
	}

	if len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "" {
		cfg.AllowedOrigins = []string{"*"}
	}

	return cfg, nil
}
