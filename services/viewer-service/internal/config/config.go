package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort string
	GrpcPort   string

	DatabaseURL string
	RedisURL    string

	JWTSecret string
	JWTExpiry time.Duration

	KafkaBrokers []string

	AllowedOrigins []string

	HistoryRetentionDays int
	WatchProgressInterval int
	MaxBookmarks         int
	MaxWatchLater        int

	LogLevel string
	Env      string
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("SERVER_PORT", "4003")
	v.SetDefault("GRPC_PORT", "9103")
	v.SetDefault("JWT_EXPIRY", "15m")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("ENV", "development")
	v.SetDefault("HISTORY_RETENTION_DAYS", 90)
	v.SetDefault("WATCH_PROGRESS_INTERVAL", 30)
	v.SetDefault("MAX_BOOKMARKS", 5000)
	v.SetDefault("MAX_WATCH_LATER", 1000)

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

		AllowedOrigins: strings.Split(v.GetString("ALLOWED_ORIGINS"), ","),

		HistoryRetentionDays:  getIntEnv(v, "HISTORY_RETENTION_DAYS", 90),
		WatchProgressInterval: getIntEnv(v, "WATCH_PROGRESS_INTERVAL", 30),
		MaxBookmarks:          getIntEnv(v, "MAX_BOOKMARKS", 5000),
		MaxWatchLater:         getIntEnv(v, "MAX_WATCH_LATER", 1000),

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

	return cfg, nil
}

func getIntEnv(v *viper.Viper, key string, defaultVal int) int {
	val := v.GetString(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
