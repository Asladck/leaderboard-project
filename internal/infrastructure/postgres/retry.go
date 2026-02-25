package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
	Timeout     time.Duration
}

func ConnectWithRetry(ctx context.Context, cfg RetryConfig, username, password, host, port, dbname, sslmode string) (*sqlx.DB, error) {

	deadline := time.Now().Add(cfg.Timeout)
	attempt := 1

	for {
		db, err := Connect(username, password, host, port, dbname, sslmode)
		if err == nil {
			return db, nil
		}

		if time.Now().After(deadline) || attempt >= cfg.MaxAttempts {
			return nil, fmt.Errorf(
				"db connection failed after %d attempts: %w",
				attempt,
				err,
			)
		}

		time.Sleep(cfg.Delay)
		attempt++
	}
}
