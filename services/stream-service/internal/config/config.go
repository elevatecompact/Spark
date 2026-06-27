package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	S3       S3Config
	RTMP     RTMPConfig
	WebRTC   WebRTCConfig
	HLS      HLSConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Host            string
	Port            int
	GracefulTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	AllowedOrigins  []string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
}

type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type RTMPConfig struct {
	ListenAddr string
	AppName    string
	Domain     string
}

type WebRTCConfig struct {
	ICEServers   []ICEServer
	ICETrickle   bool
	MaxViewers   int
	STUNServers  []string
	TLSCertPath  string
	TLSKeyPath   string
}

type ICEServer struct {
	URLs       []string `mapstructure:"urls"`
	Username   string   `mapstructure:"username"`
	Credential string   `mapstructure:"credential"`
}

type HLSConfig struct {
	SegmentDuration int
	PlaylistLength  int
	OutputDir       string
	Qualities       []QualityConfig
}

type QualityConfig struct {
	Name    string
	Width   int
	Height  int
	Bitrate int
}

type JWTConfig struct {
	Secret     string
	Issuer     string
	Audience   string
	Expiration time.Duration
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.graceful_timeout", "30s")
	v.SetDefault("server.read_timeout", "30s")
	v.SetDefault("server.write_timeout", "30s")
	v.SetDefault("server.idle_timeout", "60s")
	v.SetDefault("server.allowed_origins", []string{"*"})

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "spark")
	v.SetDefault("database.password", "spark")
	v.SetDefault("database.dbname", "spark_stream")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.topic", "stream_events")

	v.SetDefault("s3.endpoint", "")
	v.SetDefault("s3.region", "us-east-1")
	v.SetDefault("s3.bucket", "spark-recordings")
	v.SetDefault("s3.access_key", "")
	v.SetDefault("s3.secret_key", "")
	v.SetDefault("s3.use_ssl", true)

	v.SetDefault("rtmp.listen_addr", ":1935")
	v.SetDefault("rtmp.app_name", "live")
	v.SetDefault("rtmp.domain", "live.spark.com")

	v.SetDefault("webrtc.ice_trickle", true)
	v.SetDefault("webrtc.max_viewers", 1000)
	v.SetDefault("webrtc.stun_servers", []string{"stun:stun.l.google.com:19302"})

	v.SetDefault("hls.segment_duration", 4)
	v.SetDefault("hls.playlist_length", 10)
	v.SetDefault("hls.output_dir", "/data/hls")
	v.SetDefault("hls.qualities", []map[string]interface{}{
		{"name": "source", "width": 1920, "height": 1080, "bitrate": 8000},
		{"name": "720p", "width": 1280, "height": 720, "bitrate": 5000},
		{"name": "480p", "width": 854, "height": 480, "bitrate": 2500},
		{"name": "360p", "width": 640, "height": 360, "bitrate": 1200},
	})

	v.SetDefault("jwt.secret", "change-me-in-production")
	v.SetDefault("jwt.issuer", "spark")
	v.SetDefault("jwt.audience", "spark-stream")
	v.SetDefault("jwt.expiration", "24h")

	v.AutomaticEnv()
	v.SetEnvPrefix("SPARK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Host:            v.GetString("server.host"),
			Port:            v.GetInt("server.port"),
			GracefulTimeout: v.GetDuration("server.graceful_timeout"),
			ReadTimeout:     v.GetDuration("server.read_timeout"),
			WriteTimeout:    v.GetDuration("server.write_timeout"),
			IdleTimeout:     v.GetDuration("server.idle_timeout"),
			AllowedOrigins:  v.GetStringSlice("server.allowed_origins"),
		},
		Database: DatabaseConfig{
			Host:            v.GetString("database.host"),
			Port:            v.GetInt("database.port"),
			User:            v.GetString("database.user"),
			Password:        v.GetString("database.password"),
			DBName:          v.GetString("database.dbname"),
			SSLMode:         v.GetString("database.sslmode"),
			MaxOpenConns:    v.GetInt("database.max_open_conns"),
			MaxIdleConns:    v.GetInt("database.max_idle_conns"),
			ConnMaxLifetime: v.GetDuration("database.conn_max_lifetime"),
		},
		Redis: RedisConfig{
			Host:     v.GetString("redis.host"),
			Port:     v.GetInt("redis.port"),
			Password: v.GetString("redis.password"),
			DB:       v.GetInt("redis.db"),
		},
		Kafka: KafkaConfig{
			Brokers: v.GetStringSlice("kafka.brokers"),
			Topic:   v.GetString("kafka.topic"),
		},
		S3: S3Config{
			Endpoint:  v.GetString("s3.endpoint"),
			Region:    v.GetString("s3.region"),
			Bucket:    v.GetString("s3.bucket"),
			AccessKey: v.GetString("s3.access_key"),
			SecretKey: v.GetString("s3.secret_key"),
			UseSSL:    v.GetBool("s3.use_ssl"),
		},
		RTMP: RTMPConfig{
			ListenAddr: v.GetString("rtmp.listen_addr"),
			AppName:    v.GetString("rtmp.app_name"),
			Domain:     v.GetString("rtmp.domain"),
		},
		WebRTC: WebRTCConfig{
			ICETrickle:  v.GetBool("webrtc.ice_trickle"),
			MaxViewers:  v.GetInt("webrtc.max_viewers"),
			STUNServers: v.GetStringSlice("webrtc.stun_servers"),
		},
		HLS: HLSConfig{
			SegmentDuration: v.GetInt("hls.segment_duration"),
			PlaylistLength:  v.GetInt("hls.playlist_length"),
			OutputDir:       v.GetString("hls.output_dir"),
		},
		JWT: JWTConfig{
			Secret:     v.GetString("jwt.secret"),
			Issuer:     v.GetString("jwt.issuer"),
			Audience:   v.GetString("jwt.audience"),
			Expiration: v.GetDuration("jwt.expiration"),
		},
	}

	qualitiesRaw := v.Get("hls.qualities")
	if qList, ok := qualitiesRaw.([]interface{}); ok {
		for _, q := range qList {
			if qm, ok := q.(map[string]interface{}); ok {
				cfg.HLS.Qualities = append(cfg.HLS.Qualities, QualityConfig{
					Name:    qm["name"].(string),
					Width:   qm["width"].(int),
					Height:  qm["height"].(int),
					Bitrate: qm["bitrate"].(int),
				})
			}
		}
	}

	return cfg, nil
}

func (c *Config) DatabaseDSN() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password +
		"@" + c.Database.Host + ":" + fmt.Sprintf("%d", c.Database.Port) +
		"/" + c.Database.DBName + "?sslmode=" + c.Database.SSLMode
}

func (c *Config) RedisAddr() string {
	return c.Redis.Host + ":" + fmt.Sprintf("%d", c.Redis.Port)
}
