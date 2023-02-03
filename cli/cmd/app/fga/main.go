package fga

import (
	"context"
	_ "embed"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/cli/application"
)

//go:embed type-definition.json
var typeDefinitionsContent string

func NewCmd(app *application.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "fga",
	}
	update := &cobra.Command{
		Use:   "update",
		Short: "updates the openFGA type definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			err := app.Start(ctx)
			if err != nil {
				return err
			}

			_, err = app.FGAClient.SetTypeDefinitions(ctx, typeDefinitionsContent)
			if err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(update)

	return rootCmd
}
