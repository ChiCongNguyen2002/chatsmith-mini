package databases

import (
	"context"
	"time"

	"chatsmith-mini/pkg/failure"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pingAttempts = 10
	pingInterval = time.Second
)

func NewPostgres(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}

	if err := pingWithRetry(ctx, pool); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func pingWithRetry(ctx context.Context, pool *pgxpool.Pool) error {
	var err error
	for attempt := 0; attempt < pingAttempts; attempt++ {
		if err = pool.Ping(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return failure.ErrWithTrace(ctx.Err())
		case <-time.After(pingInterval):
		}
	}

	return failure.ErrWithTrace(err)
}
