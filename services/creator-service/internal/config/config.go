package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URL                     string
	MaxOpenConnections      int
	MaxIdleConnections      int
	ConnectionMaxLifetime   time.Duration
	ConnectionMaxIdleTime   time.Duration
}

type RedisConfig struct {
	URL      string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type JWTConfig struct {
	Secret string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetDefault("server.port", "8083")
	v.SetDefault("server.shutdownTimeout", 30*time.Second)

	v.SetDefault("database.maxOpenConnections", 25)
	v.SetDefault("database.maxIdleConnections", 10)
	v.SetDefault("database.connectionMaxLifetime", 5*time.Minute)
	v.SetDefault("database.connectionMaxIdleTime", 1*time.Minute)

	v.SetDefault("redis.db", 0)

	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.topic", "creator-events")

	v.AutomaticEnv()
	v.SetEnvPrefix("CREATOR")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	return &Config{
		Server: ServerConfig{
			Port:            v.GetString("server.port"),
			ShutdownTimeout: v.GetDuration("server.shutdownTimeout"),
		},
		Database: DatabaseConfig{
			URL:                   v.GetString("database.url"),
			MaxOpenConnections:    v.GetInt("database.maxOpenConnections"),
			MaxIdleConnections:    v.GetInt("database.maxIdleConnections"),
			ConnectionMaxLifetime: v.GetDuration("database.connectionMaxLifetime"),
			ConnectionMaxIdleTime: v.GetDuration("database.connectionMaxIdleTime"),
		},
		Redis: RedisConfig{
			URL:      v.GetString("redis.url"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers: v.GetStringSlice("kafka.brokers"),
			Topic:   v.GetString("kafka.topic"),
		},
		JWT: JWTConfig{
			Secret: v.GetString("jwt.secret"),
		},
	}, nil
}
