package application

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/suse-skyscraper/skyscraper/cli/config"
	"github.com/suse-skyscraper/skyscraper/test/containerpool"

	"github.com/stretchr/testify/assert"
)

func TestNewApp(t *testing.T) {
	confDir, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}

	tt := []struct {
		name        string
		configDir   string
		expectedErr bool
	}{
		{
			name:        "valid config",
			configDir:   confDir,
			expectedErr: false,
		},
		{
			name:        "missing config",
			configDir:   path.Join(confDir, "testdata/missing-config"),
			expectedErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewApp(tc.configDir)
			if tc.expectedErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestApp_Start(t *testing.T) {
	ctx := context.Background()

	pool, err := containerpool.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	conf, err := pool.ContainerizedConf(ctx)
	if err != nil {
		t.Fatal(err)
	}

	app := &App{
		Config: conf,
	}

	err = app.Start(ctx)
	assert.Nil(t, err)
	app.Shutdown(ctx)
}

func TestApp_Start_BadDB(t *testing.T) {
	ctx := context.Background()

	pool, err := containerpool.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	conf, err := pool.ContainerizedConf(ctx)
	if err != nil {
		t.Fatal(err)
	}

	app := &App{
		Config: conf,
	}
	app.Config.DB.Host = "badhost"

	err = app.Start(ctx)
	assert.NotNil(t, err)
	app.Shutdown(ctx)
}

func TestApp_Start_NoDB(t *testing.T) {
	ctx := context.Background()

	pool, err := containerpool.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	conf, err := pool.ContainerizedConf(ctx)
	if err != nil {
		t.Fatal(err)
	}

	app := &App{
		Config: conf,
	}
	app.Config.DB = config.DBConfig{}

	err = app.Start(ctx)
	assert.NotNil(t, err)
	app.Shutdown(ctx)
}

func TestApp_Start_BadNats(t *testing.T) {
	ctx := context.Background()

	pool, err := containerpool.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	conf, err := pool.ContainerizedConf(ctx)
	if err != nil {
		t.Fatal(err)
	}

	app := &App{
		Config: conf,
	}
	app.Config.Nats.URL = "badhost"

	err = app.Start(ctx)
	assert.NotNil(t, err)
	app.Shutdown(ctx)
}

func TestApp_Start_BadFGA(t *testing.T) {
	ctx := context.Background()

	pool, err := containerpool.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	conf, err := pool.ContainerizedConf(ctx)
	if err != nil {
		t.Fatal(err)
	}

	app := &App{
		Config: conf,
	}
	app.Config.FGAConfig.APIScheme = "foo"

	err = app.Start(ctx)
	assert.NotNil(t, err)
	app.Shutdown(ctx)
}
