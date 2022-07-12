package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Searches struct {
	pool *pgxpool.Pool
}

func NewSearch(pool *pgxpool.Pool) *Searches {
	return &Searches{pool: pool}
}

type Searcher interface {
	SearchCloudAccounts(ctx context.Context, input SearchCloudAccountsInput) ([]CloudAccount, error)
}
