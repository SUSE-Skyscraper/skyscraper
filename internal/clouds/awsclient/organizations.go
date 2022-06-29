package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

type OrganizationsClient struct {
	client *organizations.Client
}

func NewOrganizationsClient(cfg aws.Config) *OrganizationsClient {
	client := organizations.NewFromConfig(cfg)

	return &OrganizationsClient{
		client: client,
	}
}

func (o *OrganizationsClient) ListAccounts(ctx context.Context) ([]types.Account, error) {
	accounts, err := o.client.ListAccounts(ctx, &organizations.ListAccountsInput{})

	return accounts.Accounts, err
}

func (o *OrganizationsClient) ListTagsForAccount(ctx context.Context, accountID string) ([]types.Tag, error) {
	tags, err := o.client.ListTagsForResource(ctx, &organizations.ListTagsForResourceInput{
		ResourceId: aws.String(accountID),
	})

	return tags.Tags, err
}

func (o *OrganizationsClient) TagAccount(ctx context.Context, accountID string, tags []types.Tag) error {
	_, err := o.client.TagResource(ctx, &organizations.TagResourceInput{
		ResourceId: aws.String(accountID),
		Tags:       tags,
	})

	return err
}
