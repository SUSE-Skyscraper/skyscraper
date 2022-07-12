package cloudsync

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/clouds/awsclient"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud-sync",
		Short: "Syncs cloud tags to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := syncAWSAccounts(app)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func syncAWSAccounts(app *application.App) error {
	ctx := context.Background()

	for _, tenant := range app.Config.Clouds.AWSTenants {
		cfg, err := awsclient.NewConfig(ctx, tenant.AccessKeyID, tenant.SecretAccessKey, tenant.Region)
		if err != nil {
			return err
		}
		organizationsClient := awsclient.NewOrganizationsClient(cfg)

		err = app.Repository.CreateCloudTenant(ctx, db.CreateCloudTenantParams{
			Cloud:    "AWS",
			TenantID: tenant.MasterAccountID,
			Name:     tenant.Name,
		})
		if err != nil {
			return err
		}

		accounts, err := organizationsClient.ListAccounts(ctx)
		if err != nil {
			return err
		}

		for _, account := range accounts {
			err = organizationsClient.SyncTags(ctx, app, awsclient.SyncTagsInput{
				AccountID:   aws.ToString(account.Id),
				TenantID:    tenant.MasterAccountID,
				AccountName: aws.ToString(account.Name),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
