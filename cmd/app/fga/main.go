package fga

import (
	"context"
	_ "embed"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
)

//go:embed type-definition.json
var typeDefinitionsContent string

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-fga",
		Short: "updates the openFGA type definitions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			err := app.Start(ctx)
			if err != nil {
				return err
			}

			return app.FGAClient.SetTypeDefinitions(ctx, typeDefinitionsContent)
		},
	}

	return cmd
}
