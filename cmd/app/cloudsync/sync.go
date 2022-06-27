package cloudsync

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/jackc/pgtype"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud-sync",
		Short: "Syncs cloud tags to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, tenant := range app.Config.Clouds.AWSTenants {
				creds := credentials.NewStaticCredentialsProvider(tenant.AccessKeyID, tenant.SecretAccessKey, "")
				cfg, err := config.LoadDefaultConfig(context.TODO(),
					config.WithCredentialsProvider(creds),
				)
				if err != nil {
					return err
				}
				ctx := context.TODO()

				err = app.DB.CreateCloudTenant(ctx, db.CreateCloudTenantParams{
					Cloud:    "AWS",
					TenantID: tenant.MasterAccountID,
					Name:     tenant.Name,
				})
				if err != nil {
					return err
				}

				organizationsAPI := organizations.NewFromConfig(cfg)
				accounts, err := organizationsAPI.ListAccounts(context.TODO(), &organizations.ListAccountsInput{})
				if err != nil {
					return err
				}

				for _, account := range accounts.Accounts {
					accountTags, err := organizationsAPI.ListTagsForResource(context.TODO(), &organizations.ListTagsForResourceInput{
						ResourceId: account.Id,
					})
					if err != nil {
						return err
					}

					var tags = make(map[string]interface{})
					for _, tag := range accountTags.Tags {
						tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
					}
					json := pgtype.JSONB{}
					err = json.Set(tags)
					if err != nil {
						return err
					}

					err = app.DB.CreateCloudAccountMetadata(ctx, db.CreateCloudAccountMetadataParams{
						Cloud:       "AWS",
						TenantID:    tenant.MasterAccountID,
						AccountID:   aws.ToString(account.Id),
						Name:        aws.ToString(account.Name),
						TagsCurrent: json,
						TagsDesired: json,
					})
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	return cmd
}
