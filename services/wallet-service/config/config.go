package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"service_name"`

	Port    int    `mapstructure:"port"`

	PostgresHost     string `mapstructure:"postgres_host"`
	PostgresPort     int    `mapstructure:"postgres_port"`
	PostgresUser     string `mapstructure:"postgres_user"`
	PostgresPassword string `mapstructure:"postgres_password"`
	PostgresDB       string `mapstructure:"postgres_db"`

	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`

	KafkaBrokers []string `mapstructure:"kafka_brokers"`

	MaxBalanceCents  int64 `mapstructure:"max_balance_cents"`
	PayoutMinCents   int64 `mapstructure:"payout_minimum_cents"`

	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("service_name", "wallet-service")
	v.SetDefault("port", 4009)
	v.SetDefault("postgres_host", "localhost")
	v.SetDefault("postgres_port", 5432)
	v.SetDefault("postgres_user", "spark")
	v.SetDefault("postgres_password", "spark")
	v.SetDefault("postgres_db", "spark_wallet")
	v.SetDefault("redis_addr", "localhost:6379")
	v.SetDefault("redis_password", "")
	v.SetDefault("redis_db", 0)
	v.SetDefault("kafka_brokers", []string{"localhost:9092"})
	v.SetDefault("max_balance_cents", 100000000)
	v.SetDefault("payout_minimum_cents", 5000)
	v.SetDefault("allowed_origins", []string{"*"})

	v.AutomaticEnv()
	v.SetEnvPrefix("WALLET")

	v.BindEnv("port", "WALLET_PORT")
	v.BindEnv("postgres_host", "WALLET_POSTGRES_HOST")
	v.BindEnv("postgres_port", "WALLET_POSTGRES_PORT")
	v.BindEnv("postgres_user", "WALLET_POSTGRES_USER")
	v.BindEnv("postgres_password", "WALLET_POSTGRES_PASSWORD")
	v.BindEnv("postgres_db", "WALLET_POSTGRES_DB")
	v.BindEnv("redis_addr", "WALLET_REDIS_ADDR")
	v.BindEnv("redis_password", "WALLET_REDIS_PASSWORD")
	v.BindEnv("redis_db", "WALLET_REDIS_DB")
	v.BindEnv("kafka_brokers", "WALLET_KAFKA_BROKERS")
	v.BindEnv("max_balance_cents", "WALLET_MAX_BALANCE_CENTS")
	v.BindEnv("payout_minimum_cents", "WALLET_PAYOUT_MINIMUM_CENTS")
	v.BindEnv("allowed_origins", "WALLET_ALLOWED_ORIGINS")

	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	v.SetConfigName(fmt.Sprintf("config-%s", env))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
