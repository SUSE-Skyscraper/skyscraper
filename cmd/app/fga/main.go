package fga

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
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
	assert := &cobra.Command{
		Use:   "assert",
		Short: "Runs the openFGA assertions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			err := app.Start(ctx)
			if err != nil {
				return err
			}

			ok, err := app.FGAClient.RunAssertions(ctx, typeDefinitionsContent)
			if err != nil {
				return err
			} else if !ok {
				fmt.Println("assertions failed")
				os.Exit(1)
			}

			fmt.Println("assertions passed")
			return nil
		},
	}

	rootCmd.AddCommand(update)
	rootCmd.AddCommand(assert)

	return rootCmd
}
