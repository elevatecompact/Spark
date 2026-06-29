package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"service_name"`

	Port    int `mapstructure:"port"`
	WSPort  int `mapstructure:"ws_port"`

	PostgresHost     string `mapstructure:"postgres_host"`
	PostgresPort     int    `mapstructure:"postgres_port"`
	PostgresUser     string `mapstructure:"postgres_user"`
	PostgresPassword string `mapstructure:"postgres_password"`
	PostgresDB       string `mapstructure:"postgres_db"`

	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`

	KafkaBrokers []string `mapstructure:"kafka_brokers"`
	KafkaTopic   string   `mapstructure:"kafka_topic"`

	MaxMessageLength    int `mapstructure:"max_message_length"`
	MessageHistoryLimit int `mapstructure:"message_history_limit"`
	MessageRetentionDays int `mapstructure:"message_retention_days"`
	WSHeartbeatInterval  int `mapstructure:"ws_heartbeat_interval"`

	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("service_name", "chat-service")
	v.SetDefault("port", 4005)
	v.SetDefault("ws_port", 4006)
	v.SetDefault("postgres_host", "localhost")
	v.SetDefault("postgres_port", 5432)
	v.SetDefault("postgres_user", "spark")
	v.SetDefault("postgres_password", "spark")
	v.SetDefault("postgres_db", "spark_chat")
	v.SetDefault("redis_addr", "localhost:6379")
	v.SetDefault("redis_password", "")
	v.SetDefault("redis_db", 0)
	v.SetDefault("kafka_brokers", []string{"localhost:9092"})
	v.SetDefault("kafka_topic", "chat-events")
	v.SetDefault("max_message_length", 500)
	v.SetDefault("message_history_limit", 100)
	v.SetDefault("message_retention_days", 90)
	v.SetDefault("ws_heartbeat_interval", 30)
	v.SetDefault("allowed_origins", []string{"*"})

	v.AutomaticEnv()
	v.SetEnvPrefix("CHAT")

	v.BindEnv("port", "CHAT_PORT")
	v.BindEnv("ws_port", "CHAT_WS_PORT")
	v.BindEnv("postgres_host", "CHAT_POSTGRES_HOST")
	v.BindEnv("postgres_port", "CHAT_POSTGRES_PORT")
	v.BindEnv("postgres_user", "CHAT_POSTGRES_USER")
	v.BindEnv("postgres_password", "CHAT_POSTGRES_PASSWORD")
	v.BindEnv("postgres_db", "CHAT_POSTGRES_DB")
	v.BindEnv("redis_addr", "CHAT_REDIS_ADDR")
	v.BindEnv("redis_password", "CHAT_REDIS_PASSWORD")
	v.BindEnv("redis_db", "CHAT_REDIS_DB")
	v.BindEnv("kafka_brokers", "CHAT_KAFKA_BROKERS")
	v.BindEnv("kafka_topic", "CHAT_KAFKA_TOPIC")
	v.BindEnv("max_message_length", "CHAT_MAX_MESSAGE_LENGTH")
	v.BindEnv("message_history_limit", "CHAT_MESSAGE_HISTORY_LIMIT")
	v.BindEnv("message_retention_days", "CHAT_MESSAGE_RETENTION_DAYS")
	v.BindEnv("ws_heartbeat_interval", "CHAT_WS_HEARTBEAT_INTERVAL")
	v.BindEnv("allowed_origins", "CHAT_ALLOWED_ORIGINS")

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
