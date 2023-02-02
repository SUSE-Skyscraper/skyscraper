package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/suse-skyscraper/skyscraper/api/queue"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/plugins/aws/awsclient"
)

func newWorkerCmd(app *application.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "starts the worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// listen for change tags messages
			sub, err := app.JS.PullSubscribe("PLUGIN.AWS", "aws_worker")
			if err != nil {
				return err
			}

			// Start a goroutine to print the consumer info
			go consumerInfo(sub)

			for {
				err := fetchMessage(ctx, app, sub)
				if err != nil {
					log.Println(err)
				}
			}
		},
	}

	return cmd
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

func fetchMessage(ctx context.Context, app *application.App, sub *nats.Subscription) error {
	messages, err := sub.Fetch(1, nats.MaxWait(10*time.Second))
	if err != nil {
		return err
	}

	for _, msg := range messages {
		log.Println("Received message", string(msg.Data))

		err := msg.InProgress(nats.AckWait(10 * time.Second))
		if err != nil {
			return fmt.Errorf("in progress: %w", err)
		}

		payload := queue.PluginPayload{}
		err = json.Unmarshal(msg.Data, &payload)
		if err != nil {
			return err
		}

		err = runMessage(ctx, app, msg, payload)
		if err != nil {
			_ = msg.Nak()

			return fmt.Errorf("change tags: %w", err)
		}

		err = msg.AckSync()
		if err != nil {
			return fmt.Errorf("ack sync: %w", err)
		}
	}

	return nil
}

func runMessage(ctx context.Context, app *application.App, msg *nats.Msg, payload queue.PluginPayload) error {
	switch payload.Action {
	case queue.PluginActionTagUpdate:
		return changeTags(ctx, app, payload)
	default:
		log.Println("unknown action", payload.Action)
		return msg.Ack()
	}
}

func changeTags(ctx context.Context, app *application.App, payload queue.PluginPayload) error {
	skyscraperClient := NewSkyscraperClient(app)

	account, err := skyscraperClient.getAccount(GetAccountInput{
		AccountID: payload.ResourceID,
		TenantID:  payload.TenantID,
	})
	if err != nil {
		return err
	}

	var tags []types.Tag
	for key, value := range account.Data.Attributes.TagsDesired {
		tags = append(tags, types.Tag{
			Key:   aws.String(key),
			Value: aws.String(value),
		})
	}

	i := sort.Search(len(app.Config.Clouds.AWSTenants), func(i int) bool {
		return app.Config.Clouds.AWSTenants[i].MasterAccountID == payload.TenantID
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
	err = organizations.TagAccount(ctx, payload.ResourceID, tags)
	if err != nil {
		return err
	}

	tagsToRemove := tagsToRemove(account.Data.Attributes.TagsCurrent, account.Data.Attributes.TagsDesired)
	if tagsToRemove != nil {
		err = organizations.UntagAccount(ctx, payload.ResourceID, tagsToRemove)
		if err != nil {
			return err
		}
	}

	err = skyscraperClient.updateAccount(ctx, organizations, UpdateAccountInput{
		AccountID: account.Data.Attributes.AccountID,
		TenantID:  tenant.MasterAccountID,
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
