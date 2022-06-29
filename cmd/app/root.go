package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/cmd/app/cloudsync"
	"github.com/suse-skyscraper/skyscraper/cmd/app/migrate"
	"github.com/suse-skyscraper/skyscraper/cmd/app/server"
	"github.com/suse-skyscraper/skyscraper/cmd/app/worker"
	"github.com/suse-skyscraper/skyscraper/internal/application"
)

var (
	// Used for the application state. Cobra hasn't read the environment flag yet, so we cannot set it up.
	app = &application.App{}

	rootCmd = &cobra.Command{
		Use: "cloud-dashboard",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := app.Setup()
			if err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(server.NewCmd(app))
	rootCmd.AddCommand(cloudsync.NewCmd(app))
	rootCmd.AddCommand(migrate.NewCmd(app))
	rootCmd.AddCommand(worker.NewCmd(app))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
