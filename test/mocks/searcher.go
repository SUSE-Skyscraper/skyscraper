package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/cli/db"
)

type TestSearcher struct {
	mock.Mock
}

func (t *TestSearcher) SearchCloudAccounts(ctx context.Context, input db.SearchCloudAccountsInput) ([]db.CloudAccount, error) {
	args := t.Called(ctx, input)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

var _ db.Searcher = (*TestSearcher)(nil)
