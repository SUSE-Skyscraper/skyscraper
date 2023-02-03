package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DefaultSearcher struct {
	pool *pgxpool.Pool
}

func NewSearcher(pool *pgxpool.Pool) *DefaultSearcher {
	return &DefaultSearcher{pool: pool}
}

type Searcher interface {
	SearchCloudAccounts(ctx context.Context, input SearchCloudAccountsInput) ([]CloudAccount, error)
}
