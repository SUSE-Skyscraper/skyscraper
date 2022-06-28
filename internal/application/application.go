package application

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type App struct {
	Config Config
	DB     *db.Queries
}

func (a *App) Setup() error {
	config, err := Configuration()
	if err != nil {
		return err
	}
	a.Config = config

	connection, err := pgx.Connect(context.Background(), config.DB.GetDSN())
	if err != nil {
		return err
	}
	a.DB = db.New(connection)

	return nil
}
