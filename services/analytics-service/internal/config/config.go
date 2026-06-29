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
	ClickHouse ClickHouseConfig
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

type ClickHouseConfig struct {
	URL string
}

type RedisConfig struct {
	URL string
}

type KafkaConfig struct {
	Brokers []string
	Enabled bool
	GroupID string
}

type AppConfig struct {
	EventRetentionDays      int
	AggregationIntervalSec  int
	DashboardCacheTTL       int
	RealtimeDashboards      bool
	HistoricalAnalytics     bool
	FunnelAnalysis          bool
	ReportScheduling        bool
	AnomalyDetection        bool
}

func Load() (*Config, error) {
	cfg := &Config{
		Service: ServiceConfig{
			Port:         getEnv("ANALYTICS_PORT", "4013"),
			ReadTimeout:  getIntEnv("READ_TIMEOUT", 10),
			WriteTimeout: getIntEnv("WRITE_TIMEOUT", 30),
			IdleTimeout:  getIntEnv("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL: getEnv("ANALYTICS_POSTGRES_URL", "postgres://spark:spark@localhost:5432/spark_analytics?sslmode=disable"),
		},
		ClickHouse: ClickHouseConfig{
			URL: getEnv("ANALYTICS_CLICKHOUSE_URL", ""),
		},
		Redis: RedisConfig{
			URL: getEnv("ANALYTICS_REDIS_URL", ""),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(getEnv("ANALYTICS_KAFKA_BROKERS", "localhost:9092"), ","),
			Enabled: getBoolEnv("KAFKA_ENABLED", false),
			GroupID: getEnv("KAFKA_CONSUMER_GROUP", "analytics-events-processor"),
		},
		AppConfig: AppConfig{
			EventRetentionDays:     getIntEnv("EVENT_RETENTION_DAYS", 90),
			AggregationIntervalSec: getIntEnv("AGGREGATION_INTERVAL_SECONDS", 60),
			DashboardCacheTTL:      getIntEnv("DASHBOARD_CACHE_TTL", 30),
			RealtimeDashboards:     getBoolEnv("REALTIME_DASHBOARDS", true),
			HistoricalAnalytics:    getBoolEnv("HISTORICAL_ANALYTICS", true),
			FunnelAnalysis:         getBoolEnv("FUNNEL_ANALYSIS", true),
			ReportScheduling:       getBoolEnv("REPORT_SCHEDULING", true),
			AnomalyDetection:       getBoolEnv("ANOMALY_DETECTION", false),
		},
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("ANALYTICS_POSTGRES_URL is required")
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
