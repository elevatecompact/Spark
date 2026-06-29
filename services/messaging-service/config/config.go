package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string `mapstructure:"service_name"`

	Port   int `mapstructure:"port"`
	WSPort int `mapstructure:"ws_port"`

	PostgresHost     string `mapstructure:"postgres_host"`
	PostgresPort     int    `mapstructure:"postgres_port"`
	PostgresUser     string `mapstructure:"postgres_user"`
	PostgresPassword string `mapstructure:"postgres_password"`
	PostgresDB       string `mapstructure:"postgres_db"`

	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`

	KafkaBrokers []string `mapstructure:"kafka_brokers"`

	MaxMessageLength   int  `mapstructure:"max_message_length"`
	MaxGroupSize       int  `mapstructure:"max_group_size"`

	AllowedOrigins []string `mapstructure:"allowed_origins"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("service_name", "messaging-service")
	v.SetDefault("port", 4007)
	v.SetDefault("ws_port", 4008)
	v.SetDefault("postgres_host", "localhost")
	v.SetDefault("postgres_port", 5432)
	v.SetDefault("postgres_user", "spark")
	v.SetDefault("postgres_password", "spark")
	v.SetDefault("postgres_db", "spark_messaging")
	v.SetDefault("redis_addr", "localhost:6379")
	v.SetDefault("redis_password", "")
	v.SetDefault("redis_db", 0)
	v.SetDefault("kafka_brokers", []string{"localhost:9092"})
	v.SetDefault("max_message_length", 4000)
	v.SetDefault("max_group_size", 500)
	v.SetDefault("allowed_origins", []string{"*"})

	v.AutomaticEnv()
	v.SetEnvPrefix("MSG")

	v.BindEnv("port", "MSG_PORT")
	v.BindEnv("ws_port", "MSG_WS_PORT")
	v.BindEnv("postgres_host", "MSG_POSTGRES_HOST")
	v.BindEnv("postgres_port", "MSG_POSTGRES_PORT")
	v.BindEnv("postgres_user", "MSG_POSTGRES_USER")
	v.BindEnv("postgres_password", "MSG_POSTGRES_PASSWORD")
	v.BindEnv("postgres_db", "MSG_POSTGRES_DB")
	v.BindEnv("redis_addr", "MSG_REDIS_ADDR")
	v.BindEnv("redis_password", "MSG_REDIS_PASSWORD")
	v.BindEnv("redis_db", "MSG_REDIS_DB")
	v.BindEnv("kafka_brokers", "MSG_KAFKA_BROKERS")
	v.BindEnv("max_message_length", "MSG_MAX_MESSAGE_LENGTH")
	v.BindEnv("max_group_size", "MSG_MAX_GROUP_SIZE")
	v.BindEnv("allowed_origins", "MSG_ALLOWED_ORIGINS")

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
