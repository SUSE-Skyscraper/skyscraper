package application

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type App struct {
	Config Config
	DB     *db.Queries
	JS     nats.JetStreamContext
}

func (a *App) Setup() error {
	ctx := context.Background()

	config, err := Configuration()
	if err != nil {
		return err
	}
	a.Config = config

	poolConfig, err := pgxpool.ParseConfig(config.DB.GetDSN())
	if err != nil {
		return err
	}

	database, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return err
	}
	a.DB = db.New(database)

	nc, _ := nats.Connect(nats.DefaultURL)
	js, _ := nc.JetStream(nats.PublishAsyncMaxPending(256))
	_, err = js.AddStream(&nats.StreamConfig{
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
		return err
	}

	a.JS = js

	return nil
}
