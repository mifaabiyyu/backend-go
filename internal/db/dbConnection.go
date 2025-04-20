package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, err
	}

	// Setting connection pool
	config.MaxConns = int32(maxOpenConns)

	idleDuration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	config.MaxConnIdleTime = idleDuration

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}
