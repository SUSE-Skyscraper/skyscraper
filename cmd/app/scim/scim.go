package scim

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/apikeys"
	"github.com/suse-skyscraper/skyscraper/internal/application"
)

func NewCmd(app *application.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "scim",
	}

	generateAPIKeyCmd := &cobra.Command{
		Use: "gen-api-key",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := apikeys.GenerateRandomStringURLSafe(32)
			if err != nil {
				return err
			}

			_, err = app.DB.InsertAPIKey(context.Background(), token)
			if err != nil {
				return err
			}

			fmt.Println(token)

			return nil
		},
	}

	rootCmd.AddCommand(generateAPIKeyCmd)

	return rootCmd
}
