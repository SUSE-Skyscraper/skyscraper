package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/plugins/aws/awsclient"
)

type SkyscraperClient struct {
	client *http.Client
	app    *application.App
}

func NewSkyscraperClient(app *application.App) *SkyscraperClient {
	return &SkyscraperClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		app: app,
	}
}

type UpdateTenantInput struct {
	Cloud    string
	TenantID string
	Name     string
}

func (c *SkyscraperClient) updateTenant(input UpdateTenantInput) error {
	data := payloads.CreateOrUpdateTenantPayload{
		Data: payloads.CreateOrUpdateTenantPayloadData{
			Cloud: input.Cloud,
			Name:  input.Name,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/groups/%s/tenants/%s", c.app.Config.ServerConfig.BaseURL, "AWS", input.TenantID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.app.Config.Plugin.APIKeyID)
	req.Header.Set("X-API-Secret", c.app.Config.Plugin.APISecretKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create or update resource: %s", resp.Status)
	}

	defer resp.Body.Close()

	return nil
}

type UpdateAccountInput struct {
	AccountID   string
	TenantID    string
	AccountName string
}

func (c *SkyscraperClient) updateAccount(ctx context.Context, organizations *awsclient.OrganizationsClient, input UpdateAccountInput) error {
	tags, err := organizations.SyncTags(ctx, awsclient.SyncTagsInput{
		AccountID: input.AccountID,
		TenantID:  input.TenantID,
	})
	if err != nil {
		return err
	}

	data := payloads.CreateOrUpdateResourcePayload{
		Data: payloads.CreateOrUpdateResourcePayloadData{
			AccountName: input.AccountName,
			Active:      true,
			TagsCurrent: tags,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v1/groups/%s/tenants/%s/resources/%s", c.app.Config.ServerConfig.BaseURL, "AWS", input.TenantID, input.AccountID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.app.Config.Plugin.APIKeyID)
	req.Header.Set("X-API-Secret", c.app.Config.Plugin.APISecretKey)
	req.Close = true

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create or update resource: %s", resp.Status)
	}

	defer resp.Body.Close()

	return nil
}

type GetAccountInput struct {
	AccountID string
	TenantID  string
}

func (c *SkyscraperClient) getAccount(input GetAccountInput) (responses.CloudAccountResponse, error) {
	url := fmt.Sprintf("%s/api/v1/groups/%s/tenants/%s/resources/%s", c.app.Config.ServerConfig.BaseURL, "AWS", input.TenantID, input.AccountID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return responses.CloudAccountResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.app.Config.Plugin.APIKeyID)
	req.Header.Set("X-API-Secret", c.app.Config.Plugin.APISecretKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return responses.CloudAccountResponse{}, err
	} else if resp.StatusCode != http.StatusOK {
		return responses.CloudAccountResponse{}, fmt.Errorf("failed to create or update resource: %s", resp.Status)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return responses.CloudAccountResponse{}, err
	}

	accountResponse := responses.CloudAccountResponse{}
	err = json.Unmarshal(b, &accountResponse)
	if err != nil {
		return responses.CloudAccountResponse{}, err
	}

	return accountResponse, nil
}
