package containerpool

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/nats-io/nats.go"
	openfga "github.com/openfga/go-sdk"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/suse-skyscraper/skyscraper/cli/config"
)

type Pool struct {
	pool      *dockertest.Pool
	resources []*dockertest.Resource
}

func NewPool() (*Pool, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, err
	}

	return &Pool{
		pool: pool,
	}, nil
}

func (p *Pool) Close() {
	for _, resource := range p.resources {
		_ = p.pool.Purge(resource)
	}
}

func (p *Pool) ContainerizedConf(ctx context.Context) (config.Config, error) {
	dbConfig, err := p.SetupPostgres(ctx)
	if err != nil {
		return config.Config{}, fmt.Errorf("could not start postgresql resource: %s", err)
	}

	natsConfig, err := p.SetupNats(ctx)
	if err != nil {
		return config.Config{}, fmt.Errorf("could not start nats resource: %s", err)
	}

	fgaConfig, err := p.SetupOpenFGA(ctx)
	if err != nil {
		return config.Config{}, fmt.Errorf("could not start OpenFGA resource: %s", err)
	}

	return config.Config{
		DB:        dbConfig,
		Nats:      natsConfig,
		FGAConfig: fgaConfig,
	}, nil
}

func (p *Pool) SetupOpenFGA(_ context.Context) (config.FGAConfig, error) {
	resource, err := p.pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "openfga/openfga",
		Tag:          "latest",
		Cmd:          []string{"run"},
		ExposedPorts: []string{"8080/tcp"},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return config.FGAConfig{}, err
	}

	url := resource.GetHostPort("8080/tcp")

	err = resource.Expire(120)
	if err != nil {
		return config.FGAConfig{}, err
	}

	var storeID string

	p.pool.MaxWait = 120 * time.Second
	err = p.pool.Retry(func() error {
		configuration, err := openfga.NewConfiguration(openfga.Configuration{
			ApiScheme: "http",
			ApiHost:   url,
		})
		if err != nil {
			return nil
		}

		apiClient := openfga.NewAPIClient(configuration)

		//nolint:bodyclose
		resp, _, err := apiClient.OpenFgaApi.CreateStore(context.Background()).Body(openfga.CreateStoreRequest{
			Name: "Integration Test",
		}).Execute()
		if err != nil {
			return err
		}

		storeID = resp.GetId()

		return nil
	})
	if err != nil {
		return config.FGAConfig{}, err
	}
	p.resources = append(p.resources, resource)

	conf := config.FGAConfig{
		APIScheme: "http",
		APIHost:   url,
		StoreID:   storeID,
	}

	return conf, nil
}

func (p *Pool) SetupNats(_ context.Context) (config.NatsConfig, error) {
	resource, err := p.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "nats",
		Tag:        "latest",
		Cmd:        []string{"-js"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return config.NatsConfig{}, err
	}

	url := resource.GetHostPort("4222/tcp")

	err = resource.Expire(120)
	if err != nil {
		return config.NatsConfig{}, err
	}

	p.pool.MaxWait = 120 * time.Second
	err = p.pool.Retry(func() error {
		_, err := nats.Connect(url)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return config.NatsConfig{}, err
	}
	p.resources = append(p.resources, resource)

	conf := config.NatsConfig{
		URL: url,
	}

	return conf, nil
}

func (p *Pool) SetupPostgres(ctx context.Context) (config.DBConfig, error) {
	resource, err := p.pool.RunWithOptions(&dockertest.RunOptions{
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
		return config.DBConfig{}, err
	}

	host := resource.GetBoundIP("5432/tcp")
	port := resource.GetPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://skyscraper:secret@%s:%s/skyscraper?sslmode=disable", host, port)

	err = resource.Expire(120)
	if err != nil {
		return config.DBConfig{}, err
	}

	p.pool.MaxWait = 120 * time.Second
	err = p.pool.Retry(func() error {
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
		return config.DBConfig{}, err
	}

	portParsed, err := strconv.Atoi(port)
	if err != nil {
		return config.DBConfig{}, err
	}

	p.resources = append(p.resources, resource)

	conf := config.DBConfig{
		User:     "skyscraper",
		Password: "secret",
		Database: "skyscraper",
		Port:     int64(portParsed),
		Host:     host,
	}

	return conf, nil
}
