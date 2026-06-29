package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	WSPort     string
	GrpcPort   string

	DatabaseURL string
	RedisURL    string

	JWTSecret string
	JWTExpiry time.Duration

	KafkaBrokers []string

	AllowedOrigins []string

	MaxMessageLength     int
	MessageHistoryLimit  int
	MessageRetentionDays int
	WSHeartbeatInterval  int

	LogLevel string
	Env      string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("SERVER_PORT", "4005")
	v.SetDefault("WS_PORT", "4006")
	v.SetDefault("GRPC_PORT", "9105")
	v.SetDefault("JWT_EXPIRY", "15m")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("ENV", "development")
	v.SetDefault("MAX_MESSAGE_LENGTH", 500)
	v.SetDefault("MESSAGE_HISTORY_LIMIT", 100)
	v.SetDefault("MESSAGE_RETENTION_DAYS", 90)
	v.SetDefault("WS_HEARTBEAT_INTERVAL", 30)

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
		WSPort:     v.GetString("WS_PORT"),
		GrpcPort:   v.GetString("GRPC_PORT"),

		DatabaseURL: v.GetString("DATABASE_URL"),
		RedisURL:    v.GetString("REDIS_URL"),

		JWTSecret: v.GetString("JWT_SECRET"),
		JWTExpiry: jwtExpiry,

		KafkaBrokers: strings.Split(v.GetString("KAFKA_BROKERS"), ","),

		AllowedOrigins: strings.Split(v.GetString("ALLOWED_ORIGINS"), ","),

		MaxMessageLength:     v.GetInt("MAX_MESSAGE_LENGTH"),
		MessageHistoryLimit:  v.GetInt("MESSAGE_HISTORY_LIMIT"),
		MessageRetentionDays: v.GetInt("MESSAGE_RETENTION_DAYS"),
		WSHeartbeatInterval:  v.GetInt("WS_HEARTBEAT_INTERVAL"),

		LogLevel: v.GetString("LOG_LEVEL"),
		Env:      v.GetString("ENV"),
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = "change-me-in-production"
	}

	if len(cfg.AllowedOrigins) == 1 && cfg.AllowedOrigins[0] == "" {
		cfg.AllowedOrigins = []string{"*"}
	}

	if len(cfg.KafkaBrokers) == 1 && cfg.KafkaBrokers[0] == "" {
		cfg.KafkaBrokers = []string{}
	}

	if cfg.MaxMessageLength <= 0 {
		cfg.MaxMessageLength = 500
	}
	if cfg.MessageHistoryLimit <= 0 {
		cfg.MessageHistoryLimit = 100
	}
	if cfg.MessageRetentionDays <= 0 {
		cfg.MessageRetentionDays = 90
	}
	if cfg.WSHeartbeatInterval <= 0 {
		cfg.WSHeartbeatInterval = 30
	}

	return cfg, nil
}
