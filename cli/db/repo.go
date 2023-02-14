package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Close()
}

type Repository interface {
	Querier
	WithTx(tx pgx.Tx) Repository
}

type DefaultRepo struct {
	*Queries
}

func (q *DefaultRepo) WithTx(tx pgx.Tx) Repository {
	return &DefaultRepo{
		Queries: q.Queries.WithTx(tx),
	}
}

func NewRepo(db DBTX) *DefaultRepo {
	return &DefaultRepo{New(db)}
}
