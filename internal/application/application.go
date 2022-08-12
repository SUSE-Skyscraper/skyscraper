package application

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
	openfga "github.com/openfga/go-sdk"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/fga"
)

type App struct {
	Config       Config
	JS           nats.JetStreamContext
	Repository   db.RepositoryQueries
	natsConn     *nats.Conn
	postgresPool *pgxpool.Pool
	FGAClient    fga.Client
}

func NewApp(configDir string) (*App, error) {
	configurator := NewConfigurator(configDir)
	config, err := configurator.Parse()
	if err != nil {
		return &App{}, err
	}
	return &App{
		Config: config,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	database, pool, err := setupDatabase(ctx, a.Config)
	if err != nil {
		return err
	}
	a.Repository = db.NewRepository(pool, database)
	a.postgresPool = pool

	js, nc, err := setupNats(ctx, a.Config)
	if err != nil {
		return err
	}
	a.JS = js
	a.natsConn = nc

	apiClient, err := setupFGA(ctx, a.Config)
	if err != nil {
		return err
	}
	a.FGAClient = fga.NewClient(apiClient)

	return nil
}

func (a *App) Shutdown(_ context.Context) {
	if a.postgresPool != nil {
		a.postgresPool.Close()
	}

	if a.natsConn != nil {
		a.natsConn.Close()
	}
}

func setupDatabase(ctx context.Context, config Config) (*db.Queries, *pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(config.DB.GetDSN())
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, nil, err
	}
	database := db.New(pool)

	return database, pool, nil
}

func setupFGA(_ context.Context, config Config) (*openfga.APIClient, error) {
	configuration, err := openfga.NewConfiguration(openfga.Configuration{
		ApiScheme: config.FGAConfig.APIScheme,
		ApiHost:   config.FGAConfig.APIHost,
		StoreId:   config.FGAConfig.StoreID,
	})
	if err != nil {
		return nil, err
	}

	apiClient := openfga.NewAPIClient(configuration)

	return apiClient, nil
}

func setupNats(_ context.Context, conf Config) (nats.JetStreamContext, *nats.Conn, error) {
	nc, _ := nats.Connect(conf.Nats.URL)
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
	_, err := js.AddStream(&nats.StreamConfig{
		Name:       "TAGS",
		Subjects:   []string{"TAGS.*"},
		Storage:    nats.FileStorage,
		Retention:  nats.WorkQueuePolicy,
		Discard:    nats.DiscardNew,
		Duplicates: 5 * time.Minute,
		MaxMsgs:    -1,
		MaxBytes:   -1,
	})
	if err != nil {
		return nil, nil, err
	}

	return js, nc, nil
}
