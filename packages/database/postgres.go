package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGConfig struct {
	URL            string
	MaxConnections int
	MinConnections int
	MaxIdleTime    time.Duration
	HealthCheckInterval time.Duration
}

func DefaultPGConfig(url string) PGConfig {
	return PGConfig{
		URL:                url,
		MaxConnections:     20,
		MinConnections:     2,
		MaxIdleTime:        30 * time.Minute,
		HealthCheckInterval: 1 * time.Minute,
	}
}

func NewPool(ctx context.Context, cfg PGConfig) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxConnections)
	poolCfg.MinConns = int32(cfg.MinConnections)
	poolCfg.MaxConnLifetime = cfg.MaxIdleTime
	poolCfg.HealthCheckPeriod = cfg.HealthCheckInterval

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func MustPool(ctx context.Context, cfg PGConfig) *pgxpool.Pool {
	pool, err := NewPool(ctx, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create database pool: %v", err))
	}
	return pool
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		path := filepath.Join(migrationsDir, f)
		sql, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", f, err)
		}

		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", f, err)
		}
	}

	return nil
}

func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return pool.Ping(ctx)
}

func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}

func Transaction(ctx context.Context, pool *pgxpool.Pool, fn func(context.Context) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
