package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/jackc/pgtype"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
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

func (o *OrganizationsClient) UntagAccount(ctx context.Context, accountID string, tagKeys []string) error {
	_, err := o.client.UntagResource(ctx, &organizations.UntagResourceInput{
		ResourceId: aws.String(accountID),
		TagKeys:    tagKeys,
	})

	return err
}

type SyncTagsInput struct {
	AccountID   string
	TenantID    string
	AccountName string
}

func (o *OrganizationsClient) SyncTags(ctx context.Context, app *application.App, input SyncTagsInput) error {
	accountTags, err := o.ListTagsForAccount(ctx, input.AccountID)
	if err != nil {
		return err
	}

	var tags = make(map[string]string)
	for _, tag := range accountTags {
		tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
	}
	json := pgtype.JSONB{}
	err = json.Set(tags)
	if err != nil {
		return err
	}

	err = app.DB.CreateCloudAccount(ctx, db.CreateCloudAccountParams{
		Cloud:       "AWS",
		TenantID:    input.TenantID,
		AccountID:   input.AccountID,
		Name:        input.AccountName,
		TagsCurrent: json,
		TagsDesired: json,
	})
	if err != nil {
		return err
	}

	return nil
}
