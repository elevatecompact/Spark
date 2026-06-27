package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

type TLSConfig struct {
	CertFile   string
	KeyFile    string
	CAFile     string
	ServerName string
	SkipVerify bool
}

func LoadTLS(cfg TLSConfig) (*tls.Config, error) {
	tlsCfg := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	if cfg.SkipVerify {
		tlsCfg.InsecureSkipVerify = true
		return tlsCfg, nil
	}

	if cfg.CertFile != "" && cfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS cert/key pair: %w", err)
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	if cfg.CAFile != "" {
		caCert, err := os.ReadFile(cfg.CAFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert: %w", err)
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		tlsCfg.RootCAs = caPool
	}

	if cfg.ServerName != "" {
		tlsCfg.ServerName = cfg.ServerName
	}

	return tlsCfg, nil
}

func MustLoadTLS(cfg TLSConfig) *tls.Config {
	tlsCfg, err := LoadTLS(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to load TLS config: %v", err))
	}
	return tlsCfg
}
