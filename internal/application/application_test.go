package application

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApp_Start(t *testing.T) {
	ctx := context.Background()
	app := App{}
	configDir, err := filepath.Abs("../..")
	assert.Nil(t, err)

	err = app.Start(ctx, configDir)
	assert.Nil(t, err)

	err = app.postgresPool.Ping(ctx)
	assert.Nil(t, err)
}

func TestApp_Shutdown(t *testing.T) {
	ctx := context.Background()
	app := App{}
	configDir, err := filepath.Abs("../..")
	assert.Nil(t, err)

	err = app.Start(ctx, configDir)
	assert.Nil(t, err)

	app.Shutdown(ctx)
}
