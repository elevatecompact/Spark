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
	Stripe    StripeConfig
	PayPal    PayPalConfig
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

type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	PublishableKey string
}

type PayPalConfig struct {
	ClientID     string
	ClientSecret string
	WebhookID    string
}

type AppConfig struct {
	StripeEnabled    bool
	PayPalEnabled    bool
	SaveMethods      bool
	RefundsEnabled   bool
	DisputesEnabled  bool
	DefaultCurrency  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Service: ServiceConfig{
			Port:         getEnv("PAYMENT_PORT", "4012"),
			ReadTimeout:  getIntEnv("READ_TIMEOUT", 10),
			WriteTimeout: getIntEnv("WRITE_TIMEOUT", 30),
			IdleTimeout:  getIntEnv("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL: getEnv("PAYMENT_DB_URL", "postgres://spark:spark@localhost:5432/spark_payment?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("PAYMENT_REDIS_URL", ""),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Enabled: getBoolEnv("KAFKA_ENABLED", false),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", "sk_test_noop"),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", "whsec_noop"),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", "pk_test_noop"),
		},
		PayPal: PayPalConfig{
			ClientID:     getEnv("PAYPAL_CLIENT_ID", ""),
			ClientSecret: getEnv("PAYPAL_CLIENT_SECRET", ""),
			WebhookID:    getEnv("PAYPAL_WEBHOOK_ID", ""),
		},
		AppConfig: AppConfig{
			StripeEnabled:   getBoolEnv("STRIPE_ENABLED", true),
			PayPalEnabled:   getBoolEnv("PAYPAL_ENABLED", false),
			SaveMethods:     getBoolEnv("SAVING_PAYMENT_METHODS", true),
			RefundsEnabled:  getBoolEnv("REFUNDS_ENABLED", true),
			DisputesEnabled: getBoolEnv("DISPUTES_ENABLED", true),
			DefaultCurrency: getEnv("DEFAULT_CURRENCY", "USD"),
		},
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("PAYMENT_DB_URL is required")
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
