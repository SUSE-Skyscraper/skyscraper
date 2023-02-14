package scim

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/cli/db"

	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth/apikeys"
)

func NewCmd(app *application.App) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "scim",
	}

	generateAPIKeyCmd := &cobra.Command{
		Use: "gen-api-key",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiKeyGenerator := apikeys.NewGenerator(app)
			hash, token, err := apiKeyGenerator.Generate()
			if err != nil {
				return err
			}

			ctx := context.Background()

			tx, err := app.PostgresPool.Begin(ctx)
			if err != nil {
				return err
			}

			defer func(tx pgx.Tx, ctx context.Context) {
				_ = tx.Rollback(ctx)
			}(tx, ctx)

			repo := app.Repo.WithTx(tx)

			err = repo.DeleteScimAPIKey(ctx)
			if err != nil {
				return err
			}

			apiKey, err := repo.InsertAPIKey(ctx, db.InsertAPIKeyParams{
				Encodedhash: hash,
				System:      true,
				Owner:       "SCIM",
				Description: sql.NullString{String: "SCIM API key", Valid: true},
			})
			if err != nil {
				return err
			}

			_, err = repo.InsertScimAPIKey(ctx, apiKey.ID)
			if err != nil {
				return err
			}

			err = tx.Commit(ctx)
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
