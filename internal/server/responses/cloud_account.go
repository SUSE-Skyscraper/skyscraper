package responses

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type CloudAccountItemAttributes struct {
	CloudProvider     string      `json:"cloud_provider"`
	TenantID          string      `json:"tenant_id"`
	AccountID         string      `json:"account_id"`
	Name              string      `json:"name"`
	Active            bool        `json:"active"`
	TagsCurrent       interface{} `json:"tags_current"`
	TagsDesired       interface{} `json:"tags_desired"`
	TagsDriftDetected bool        `json:"tags_drift_detected"`
	CreatedAt         string      `json:"created_at"`
	UpdatedAt         string      `json:"updated_at"`
}

type CloudAccountItem struct {
	ID         string                     `json:"id"`
	Type       ObjectResponseType         `json:"type"`
	Attributes CloudAccountItemAttributes `json:"attributes"`
}

func newCloudAccount(account db.CloudAccount) CloudAccountItem {
	return CloudAccountItem{
		ID:   account.ID.String(),
		Type: ObjectResponseTypeCloudAccount,
		Attributes: CloudAccountItemAttributes{
			CloudProvider:     account.Cloud,
			TenantID:          account.TenantID,
			AccountID:         account.AccountID,
			Name:              account.Name,
			Active:            account.Active,
			TagsCurrent:       account.TagsCurrent.Get(),
			TagsDesired:       account.TagsDesired.Get(),
			TagsDriftDetected: account.TagsDriftDetected,
			CreatedAt:         account.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         account.UpdatedAt.Format(time.RFC3339),
		},
	}
}

type CloudAccountResponse struct {
	Data CloudAccountItem `json:"data"`
}

func (rd *CloudAccountResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewCloudAccountResponse(account db.CloudAccount) *CloudAccountResponse {
	cloudAccount := newCloudAccount(account)
	return &CloudAccountResponse{
		Data: cloudAccount,
	}
}

type CloudAccountListResponse struct {
	Data []CloudAccountItem `json:"data"`
}

func (rd *CloudAccountListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewCloudAccountListResponse(accounts []db.CloudAccount) *CloudAccountListResponse {
	list := make([]CloudAccountItem, len(accounts))
	for i, account := range accounts {
		list[i] = newCloudAccount(account)
	}

	return &CloudAccountListResponse{
		Data: list,
	}
}
