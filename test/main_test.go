package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	openfga "github.com/openfga/go-sdk"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/suse-skyscraper/skyscraper/internal/application"
)

var app *application.App

func TestMain(m *testing.M) {
	ctx := context.Background()
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	postgresResource, dbConfig, err := setupPostgres(ctx, pool)
	if err != nil {
		log.Fatalf("could not start postgresql resource: %s", err)
	}

	natsResource, natsConfig, err := setupNats(ctx, pool)
	if err != nil {
		log.Fatalf("could not start nats resource: %s", err)
	}

	fgaResource, fgaConfig, err := setupOpenFGA(ctx, pool)
	if err != nil {
		log.Fatalf("could not start openfga resource: %s", err)
	}

	app = &application.App{
		Config: application.Config{
			DB:        dbConfig,
			Nats:      natsConfig,
			FGAConfig: fgaConfig,
		},
	}
	err = app.Start(ctx)
	if err != nil {
		log.Fatalf("could not start application: %s", err)
	}

	code := m.Run()

	app.Shutdown(ctx)

	if err := pool.Purge(postgresResource); err != nil {
		log.Fatalf("could not purge postgres resource: %s", err)
	}

	if err := pool.Purge(natsResource); err != nil {
		log.Fatalf("could not purge nats resource: %s", err)
	}

	if err := pool.Purge(fgaResource); err != nil {
		log.Fatalf("could not purge OpenFGA resource: %s", err)
	}

	os.Exit(code)
}

func setupOpenFGA(_ context.Context, pool *dockertest.Pool) (*dockertest.Resource, application.FGAConfig, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "openfga/openfga",
		Tag:          "latest",
		Cmd:          []string{"run"},
		ExposedPorts: []string{"8080/tcp"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, application.FGAConfig{}, err
	}

	url := resource.GetHostPort("8080/tcp")

	err = resource.Expire(120)
	if err != nil {
		return nil, application.FGAConfig{}, err
	}

	var storeID string

	pool.MaxWait = 120 * time.Second
	err = pool.Retry(func() error {
		configuration, err := openfga.NewConfiguration(openfga.Configuration{
			ApiScheme: "http",
			ApiHost:   url,
		})
		if err != nil {
			return nil
		}

		apiClient := openfga.NewAPIClient(configuration)

		resp, _, err := apiClient.OpenFgaApi.CreateStore(context.Background()).Body(openfga.CreateStoreRequest{
			Name: openfga.PtrString("Integration Test"),
		}).Execute()
		if err != nil {
			return err
		}

		storeID = resp.GetId()

		return nil
	})
	if err != nil {
		return nil, application.FGAConfig{}, err
	}

	config := application.FGAConfig{
		APIScheme: "http",
		APIHost:   url,
		StoreID:   storeID,
	}

	return resource, config, nil
}

func setupNats(_ context.Context, pool *dockertest.Pool) (*dockertest.Resource, application.NatsConfig, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "nats",
		Tag:        "latest",
		Cmd:        []string{"-js"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, application.NatsConfig{}, err
	}

	url := resource.GetHostPort("4222/tcp")

	err = resource.Expire(120)
	if err != nil {
		return nil, application.NatsConfig{}, err
	}

	pool.MaxWait = 120 * time.Second
	err = pool.Retry(func() error {
		_, err := nats.Connect(url)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, application.NatsConfig{}, err
	}

	config := application.NatsConfig{
		URL: url,
	}

	return resource, config, nil
}

func setupPostgres(ctx context.Context, pool *dockertest.Pool) (*dockertest.Resource, application.DBConfig, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=skyscraper",
			"POSTGRES_DB=skyscraper",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, application.DBConfig{}, err
	}

	host := resource.GetBoundIP("5432/tcp")
	port := resource.GetPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://skyscraper:secret@%s:%s/skyscraper?sslmode=disable", host, port)

	err = resource.Expire(120)
	if err != nil {
		return nil, application.DBConfig{}, err
	}

	pool.MaxWait = 120 * time.Second
	err = pool.Retry(func() error {
		conn, err := pgx.Connect(ctx, databaseURL)
		if err != nil {
			return err
		}

		err = conn.Ping(ctx)
		if err != nil {
			return err
		}

		return conn.Close(ctx)
	})
	if err != nil {
		return nil, application.DBConfig{}, err
	}

	portParsed, err := strconv.Atoi(port)
	if err != nil {
		return nil, application.DBConfig{}, err
	}

	config := application.DBConfig{
		User:     "skyscraper",
		Password: "secret",
		Database: "skyscraper",
		Port:     int64(portParsed),
		Host:     host,
	}

	return resource, config, nil
}
