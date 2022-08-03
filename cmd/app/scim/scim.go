package scim

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/auth/apikeys"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

func NewCmd(app *application.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "scim",
	}

	generateAPIKeyCmd := &cobra.Command{
		Use: "gen-api-key",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKeyGenerator := apikeys.NewGenerator(app.Config.Argon2Config.MemoryCost, app.Config.Argon2Config.TimeCost, app.Config.Argon2Config.Parallelism)
			hash, token, err := apiKeyGenerator.Generate()
			if err != nil {
				return err
			}

			ctx := context.Background()

			repo, err := app.Repository.Begin(ctx)
			if err != nil {
				return err
			}

			defer func(repo db.RepositoryQueries, ctx context.Context) {
				_ = repo.Rollback(ctx)
			}(repo, ctx)

			err = repo.DeleteScimAPIKey(ctx)
			if err != nil {
				return err
			}

			_, err = repo.InsertScimAPIKey(context.Background(), hash)
			if err != nil {
				return err
			}

			err = repo.Commit(ctx)
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
