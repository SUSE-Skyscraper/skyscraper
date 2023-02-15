package helpers

import (
	"context"
	"log"

	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/test/containerpool"
)

type AppContainerized struct {
	App           *application.App
	containerPool *containerpool.Pool
}

func (t *AppContainerized) Close() {
	t.containerPool.Close()
}

func NewContainerizedApp(ctx context.Context) (*AppContainerized, error) {
	pool, err := containerpool.NewPool()
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	conf, err := pool.ContainerizedConf(ctx)
	if err != nil {
		log.Fatalf("Could not setup containerized conf: %s", err)
	}

	app := &application.App{
		Config: conf,
	}

	return &AppContainerized{
		App:           app,
		containerPool: pool,
	}, nil
}
