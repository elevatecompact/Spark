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
	Push      PushConfig
	Email     EmailConfig
	SMS       SMSConfig
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

type PushConfig struct {
	FCMKey      string
	APNSKeyID   string
	APNSTeamID  string
	APNSKeyPath string
}

type EmailConfig struct {
	SendGridAPIKey string
}

type SMSConfig struct {
	TwilioSID      string
	TwilioToken    string
	TwilioPhone    string
}

type AppConfig struct {
	PushEnabled    bool
	EmailEnabled   bool
	SMSEnabled     bool
	InAppEnabled   bool
	DigestsEnabled bool
}

func Load() (*Config, error) {
	cfg := &Config{
		Service: ServiceConfig{
			Port:         getEnv("NOTIFICATION_PORT", "4014"),
			ReadTimeout:  getIntEnv("READ_TIMEOUT", 10),
			WriteTimeout: getIntEnv("WRITE_TIMEOUT", 30),
			IdleTimeout:  getIntEnv("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL: getEnv("NOTIFICATION_DB_URL", "postgres://spark:spark@localhost:5432/spark_notification?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("NOTIFICATION_REDIS_URL", ""),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("NOTIFICATION_KAFKA_BROKERS", "localhost:9092"), ","),
			Enabled: getBoolEnv("KAFKA_ENABLED", false),
		},
		Push: PushConfig{
			FCMKey:      getEnv("PUSH_FCM_SERVER_KEY", ""),
			APNSKeyID:   getEnv("PUSH_APNS_KEY_ID", ""),
			APNSTeamID:  getEnv("PUSH_APNS_TEAM_ID", ""),
			APNSKeyPath: getEnv("PUSH_APNS_KEY_PATH", ""),
		},
		Email: EmailConfig{
			SendGridAPIKey: getEnv("SENDGRID_API_KEY", ""),
		},
		SMS: SMSConfig{
			TwilioSID:   getEnv("TWILIO_ACCOUNT_SID", ""),
			TwilioToken: getEnv("TWILIO_AUTH_TOKEN", ""),
			TwilioPhone: getEnv("TWILIO_PHONE_NUMBER", ""),
		},
		AppConfig: AppConfig{
			PushEnabled:    getBoolEnv("PUSH_ENABLED", true),
			EmailEnabled:   getBoolEnv("EMAIL_ENABLED", true),
			SMSEnabled:     getBoolEnv("SMS_ENABLED", false),
			InAppEnabled:   getBoolEnv("INAPP_ENABLED", true),
			DigestsEnabled: getBoolEnv("DIGESTS_ENABLED", true),
		},
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("NOTIFICATION_DB_URL is required")
	}
	return nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
func getIntEnv(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}
func getBoolEnv(key string, defaultVal bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return defaultVal
}
