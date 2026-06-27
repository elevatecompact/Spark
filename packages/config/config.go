package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName string
	Environment string
	LogLevel    string

	Server struct {
		Port     int
		GrpcPort int
	}

	Database struct {
		URL            string
		MaxConnections int
		MinConnections int
	}

	Redis struct {
		URL      string
		Password string
	}

	Kafka struct {
		Brokers  []string
		ClientID string
	}

	JWT struct {
		Secret   string
		Expiry   string
		Issuer   string
		Audience string
	}

	Observability struct {
		OTLEndpoint string
		SentryDSN   string
	}
}

func Load(serviceName string) (*Config, error) {
	v := viper.New()
	v.SetConfigName(fmt.Sprintf("%s-config", serviceName))
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")

	v.SetEnvPrefix("SPARK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("environment", "development")
	v.SetDefault("log_level", "info")
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.max_connections", 20)
	v.SetDefault("database.min_connections", 2)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	cfg.ServiceName = serviceName
	return &cfg, nil
}
