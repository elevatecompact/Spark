package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/elevatecompact/spark/services/gift-service/internal/service"
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
	CardCodeLength     int
	CardExpiryDays     int
	MaxGiftAmountCents int64
	MinGiftAmountCents int64
	CampaignMaxDays    int
	GiftSendingEnabled bool
	GiftCardsEnabled    bool
	CampaignMatching    bool
	AnonymousGifting    bool
	GiftLeaderboard     bool
	RateLimitGifts      int
	RateLimitCards      int
	RateLimitCampaigns  int
}

func Load() (*Config, error) {
	cfg := &Config{
		Service: ServiceConfig{
			Port:         getEnv("GIFT_PORT", "4011"),
			ReadTimeout:  getIntEnv("READ_TIMEOUT", 10),
			WriteTimeout: getIntEnv("WRITE_TIMEOUT", 30),
			IdleTimeout:  getIntEnv("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL: getEnv("GIFT_DB_URL", "postgres://spark:spark@localhost:5432/spark_gift?sslmode=disable"),
		},
		Redis: RedisConfig{
			URL: getEnv("GIFT_REDIS_URL", ""),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Enabled: getBoolEnv("KAFKA_ENABLED", false),
		},
		AppConfig: AppConfig{
			CardCodeLength:     getIntEnv("GIFT_CARD_CODE_LENGTH", 12),
			CardExpiryDays:     getIntEnv("GIFT_CARD_EXPIRY_DAYS", 365),
			MaxGiftAmountCents: int64(getIntEnv("MAX_GIFT_AMOUNT_CENTS", 1000000)),
			MinGiftAmountCents: int64(getIntEnv("MIN_GIFT_AMOUNT_CENTS", 100)),
			CampaignMaxDays:    getIntEnv("CAMPAIGN_MAX_DURATION_DAYS", 30),
			GiftSendingEnabled: getBoolEnv("GIFT_SENDING_ENABLED", true),
			GiftCardsEnabled:    getBoolEnv("GIFT_CARDS_ENABLED", true),
			CampaignMatching:    getBoolEnv("CAMPAIGN_MATCHING", true),
			AnonymousGifting:    getBoolEnv("ANONYMOUS_GIFTING", true),
			GiftLeaderboard:     getBoolEnv("GIFT_LEADERBOARD", true),
			RateLimitGifts:      getIntEnv("RATE_LIMIT_GIFTS", 50),
			RateLimitCards:      getIntEnv("RATE_LIMIT_CARDS", 5),
			RateLimitCampaigns:  getIntEnv("RATE_LIMIT_CAMPAIGNS", 3),
		},
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) ToGiftServiceConfig() service.GiftServiceConfig {
	return service.GiftServiceConfig{
		CardCodeLength:     c.AppConfig.CardCodeLength,
		CardExpiryDays:     c.AppConfig.CardExpiryDays,
		MaxGiftAmountCents: c.AppConfig.MaxGiftAmountCents,
		MinGiftAmountCents: c.AppConfig.MinGiftAmountCents,
		CampaignMaxDays:    c.AppConfig.CampaignMaxDays,
		GiftSendingEnabled:   c.AppConfig.GiftSendingEnabled,
		GiftCardsEnabled:      c.AppConfig.GiftCardsEnabled,
		CampaignMatching:      c.AppConfig.CampaignMatching,
		AnonymousGifting:      c.AppConfig.AnonymousGifting,
		GiftLeaderboard:       c.AppConfig.GiftLeaderboard,
		RateLimitGifts:      c.AppConfig.RateLimitGifts,
		RateLimitCards:      c.AppConfig.RateLimitCards,
		RateLimitCampaigns:  c.AppConfig.RateLimitCampaigns,
	}
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("GIFT_DB_URL is required")
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
