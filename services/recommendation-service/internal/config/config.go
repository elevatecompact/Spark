package config

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port        string
	DatabaseURL string
	KafkaBrokers []string
	LogLevel    string
}

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8087"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://spark:spark@localhost:5432/spark_recommendation?sslmode=disable"),
		KafkaBrokers: splitEnv(getEnv("KAFKA_BROKERS", "localhost:9092")),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	log.Info().Str("port", cfg.Port).Msg("config loaded")
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func splitEnv(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}
