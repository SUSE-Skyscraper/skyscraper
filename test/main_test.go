package test

import (
	"context"
	"os"
	"testing"

	"github.com/suse-skyscraper/skyscraper/test/helpers"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	containerizedApp, err := helpers.NewContainerizedApp(ctx)
	if err != nil {
		panic(err)
	}

	app := containerizedApp.App

	code := m.Run()
	app.Shutdown(ctx)
	containerizedApp.Close()

	os.Exit(code)
}
