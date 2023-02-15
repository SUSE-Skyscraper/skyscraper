package main

import (
	"context"
	"fmt"
	"os"

	"github.com/suse-skyscraper/skyscraper/cli/config"

	"github.com/suse-skyscraper/skyscraper/cli/application"

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
		Use: "skyscraper-aws",
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

	rootCmd.AddCommand(newSyncCmd(app))
	rootCmd.AddCommand(newWorkerCmd(app))

	err = rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
