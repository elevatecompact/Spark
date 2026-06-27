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
	"github.com/rs/zerolog/log"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, dsn string, maxOpen, maxIdle int, maxLifetime time.Duration) (*PostgresDB, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pool config: %w", err)
	}

	cfg.MaxConns = int32(maxOpen)
	cfg.MinConns = int32(maxIdle)
	cfg.MaxConnLifetime = maxLifetime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Info().Msg("PostgreSQL connection established")
	return &PostgresDB{Pool: pool}, nil
}

func (db *PostgresDB) Close() {
	db.Pool.Close()
	log.Info().Msg("PostgreSQL connection closed")
}

func (db *PostgresDB) RunMigrations(ctx context.Context, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
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
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		log.Info().Str("migration", f).Msg("Running migration")

		_, err = db.Pool.Exec(ctx, string(content))
		if err != nil {
			return fmt.Errorf("execute migration %s: %w", f, err)
		}

		log.Info().Str("migration", f).Msg("Migration completed")
	}

	return nil
}

func (db *PostgresDB) Exec(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	tag, err := db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func (db *PostgresDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgxpool.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

func (db *PostgresDB) Query(ctx context.Context, sql string, args ...interface{}) (pgxpool.Rows, error) {
	return db.Pool.Query(ctx, sql, args...)
}
