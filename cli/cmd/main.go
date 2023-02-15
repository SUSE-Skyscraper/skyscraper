package main

import (
	"context"
	"fmt"
	"os"

	"github.com/suse-skyscraper/skyscraper/cli/config"

	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/cmd/app/fga"
	"github.com/suse-skyscraper/skyscraper/cli/cmd/app/migrate"
	"github.com/suse-skyscraper/skyscraper/cli/cmd/app/scim"
	"github.com/suse-skyscraper/skyscraper/cli/cmd/app/server"

	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	app, err := application.NewApp(config.DefaultConfigDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create application: %v\n", err)
		os.Exit(1)
	}

	rootCmd := &cobra.Command{
		Use: "skyscraper",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := app.Start(ctx)
			if err != nil {
				return err
			}

			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			app.Shutdown(ctx)
		},
	}

	rootCmd.AddCommand(server.NewCmd(app))
	rootCmd.AddCommand(migrate.NewCmd(app))
	rootCmd.AddCommand(scim.NewCmd(app))
	rootCmd.AddCommand(fga.NewCmd(app))

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
