package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Service   ServiceConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	Kafka     KafkaConfig
	AppConfig AppConfig
}

type ServiceConfig struct {
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type KafkaConfig struct {
	Brokers []string
	Enabled bool
}

type AppConfig struct {
	GraceDays  int
	MaxActive  int
	TrialDays  int
}

func Load() (*Config, error) {
	cfg := &Config{
		Service: ServiceConfig{
			Port:         getEnv("PORT", "8090"),
			ReadTimeout:  getIntEnv("READ_TIMEOUT", 10),
			WriteTimeout: getIntEnv("WRITE_TIMEOUT", 30),
			IdleTimeout:  getIntEnv("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://spark:spark@localhost:5432/spark_subscription?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", ""),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Enabled: getBoolEnv("KAFKA_ENABLED", false),
		},
		AppConfig: AppConfig{
			GraceDays: getIntEnv("GRACE_DAYS", 3),
			MaxActive: getIntEnv("MAX_ACTIVE_SUBSCRIPTIONS", 50),
			TrialDays: getIntEnv("TRIAL_DAYS", 0),
		},
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getIntEnv(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return defaultVal
}

func getBoolEnv(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}
