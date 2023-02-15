package helpers

import (
	"github.com/suse-skyscraper/skyscraper/cli/config"
	"github.com/suse-skyscraper/skyscraper/test/mocks"

	"github.com/pashagolub/pgxmock"

	"github.com/suse-skyscraper/skyscraper/cli/application"
)

type AppMocked struct {
	App          *application.App
	JS           *mocks.TestJS
	Repo         *mocks.TestRepo
	Searcher     *mocks.TestSearcher
	FGAClient    *mocks.TestFGAAuthorizer
	PostgresPool pgxmock.PgxPoolIface
}

func (t *AppMocked) Close() {
	t.PostgresPool.Close()
}

func NewMockedApp() (*AppMocked, error) {
	js := new(mocks.TestJS)
	fgaClient := new(mocks.TestFGAAuthorizer)
	repo := new(mocks.TestRepo)
	searcher := new(mocks.TestSearcher)
	pool, err := pgxmock.NewPool()
	if err != nil {
		return nil, err
	}

	app := &application.App{
		Config:       config.Config{},
		JS:           js,
		FGAClient:    fgaClient,
		Repo:         repo,
		PostgresPool: pool,
		Searcher:     searcher,
	}

	return &AppMocked{
		App:          app,
		JS:           js,
		FGAClient:    fgaClient,
		PostgresPool: pool,
		Repo:         repo,
		Searcher:     searcher,
	}, nil
}
