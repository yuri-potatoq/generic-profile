package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func NewDb(ctx context.Context) (*sqlx.DB, error) {
	urlExample := "postgres://furry-profile:pass@172.30.42.187:5432/furry-profile?sslmode=disable"
	pool, err := sqlx.ConnectContext(ctx, "postgres", urlExample)
	if err != nil {
		return nil, err
	}

	if err := pool.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("fail to ping DB connection: %w", err)
	}
	return pool, nil
}
