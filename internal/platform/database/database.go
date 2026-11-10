package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(url string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
