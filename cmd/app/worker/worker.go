package worker

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/clouds/awsclient"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/workers"
)

func NewCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "starts the worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			err := startChangeTagsWorker(ctx, app)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func startChangeTagsWorker(ctx context.Context, app *application.App) error {
	// listen for change tags messages
	sub, err := app.JS.PullSubscribe("TAGS.change", "TAGS")
	if err != nil {
		return err
	}

	// Start a goroutine to print the consumer info
	go consumerInfo(sub)

	for {
		messages, err := sub.Fetch(1, nats.MaxWait(10*time.Second))
		if err != nil {
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}

			log.Println("Fetch failed", err)
			continue
		}

		for _, msg := range messages {
			log.Println("Received message", string(msg.Data))

			err := msg.InProgress(nats.AckWait(10 * time.Second))
			if err != nil {
				log.Println("InProgress", err)
				continue
			}

			err = changeTags(ctx, app, msg)
			if err != nil {
				log.Println("ChangeTags", err)
				err = msg.Nak()
				if err != nil {
					log.Println("Nak", err)
				}

				continue
			}

			err = msg.AckSync()
			if err != nil {
				log.Println("AckSync", err)
				continue
			}
		}
	}
}

func consumerInfo(sub *nats.Subscription) {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		con, err := sub.ConsumerInfo()
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("NumPending=%v NumWaiting=%v\n", con.NumPending, con.NumWaiting)
	}
}

func changeTags(ctx context.Context, app *application.App, msg *nats.Msg) error {
	payload := workers.ChangeTagsPayload{}
	err := json.Unmarshal(msg.Data, &payload)
	if err != nil {
		return err
	}

	account, err := app.Repository.FindCloudAccount(ctx, db.FindCloudAccountInput{
		ID: payload.ID,
	})
	if err != nil {
		return err
	}

	var currentTags map[string]string
	err = json.Unmarshal(account.TagsCurrent.Bytes, &currentTags)
	if err != nil {
		return err
	}

	var desiredTags map[string]string
	err = json.Unmarshal(account.TagsDesired.Bytes, &desiredTags)
	if err != nil {
		return err
	}

	if account.Cloud == "AWS" {
		err = changeAwsTags(ctx, app, changeAWSTagsInput{
			tenantID:    account.TenantID,
			accountID:   account.AccountID,
			desiredTags: desiredTags,
			currentTags: currentTags,
			accountName: payload.AccountName,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type changeAWSTagsInput struct {
	tenantID    string
	accountID   string
	accountName string
	desiredTags map[string]string
	currentTags map[string]string
}

func changeAwsTags(ctx context.Context, app *application.App, input changeAWSTagsInput) error {
	var tags []types.Tag
	for key, value := range input.desiredTags {
		tags = append(tags, types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	i := sort.Search(len(app.Config.Clouds.AWSTenants), func(i int) bool {
		return app.Config.Clouds.AWSTenants[i].MasterAccountID == input.tenantID
	})
	if i < 0 {
		return errors.New("tenant not found")
	}
	tenant := app.Config.Clouds.AWSTenants[i]

	client, err := awsclient.NewConfig(ctx, tenant.AccessKeyID, tenant.SecretAccessKey, tenant.Region)
	if err != nil {
		return err
	}

	organizations := awsclient.NewOrganizationsClient(client)
	err = organizations.TagAccount(ctx, input.accountID, tags)
	if err != nil {
		return err
	}

	tagsToRemove := tagsToRemove(input.currentTags, input.desiredTags)
	if tagsToRemove != nil {
		err = organizations.UntagAccount(ctx, input.accountID, tagsToRemove)
		if err != nil {
			return err
		}
	}

	err = organizations.SyncTags(ctx, app, awsclient.SyncTagsInput{
		TenantID:    input.tenantID,
		AccountID:   input.accountID,
		AccountName: input.accountName,
	})
	if err != nil {
		return err
	}

	return nil
}

func tagsToRemove(tagsCurrent, tagsDesired map[string]string) []string {
	var tagsToRemove []string
	for key := range tagsCurrent {
		if _, ok := tagsDesired[key]; !ok {
			tagsToRemove = append(tagsToRemove, key)
		}
	}
	return tagsToRemove
}
