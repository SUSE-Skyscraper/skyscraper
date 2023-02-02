package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/plugins/aws/awsclient"
)

func newSyncCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Syncs aws tags to the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			skyscraperClient := NewSkyscraperClient(app)

			for _, tenant := range app.Config.Clouds.AWSTenants {
				cfg, err := awsclient.NewConfig(ctx, tenant.AccessKeyID, tenant.SecretAccessKey, tenant.Region)
				if err != nil {
					return err
				}
				organizationsClient := awsclient.NewOrganizationsClient(cfg)

				err = skyscraperClient.updateTenant(UpdateTenantInput{
					TenantID: tenant.MasterAccountID,
					Cloud:    "AWS",
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
					err2 := skyscraperClient.updateAccount(ctx, organizationsClient, UpdateAccountInput{
						AccountID:   aws.ToString(account.Id),
						TenantID:    tenant.MasterAccountID,
						AccountName: aws.ToString(account.Name),
					})
					if err2 != nil {
						return err2
					}
				}
			}

			return nil
		},
	}

	return cmd
}
