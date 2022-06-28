package responses

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type CloudAccountResponse struct {
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

func (rd *CloudAccountResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
func NewCloudAccountResponse(account db.CloudAccount) *CloudAccountResponse {
	return &CloudAccountResponse{
		CloudProvider:     account.Cloud,
		TenantID:          account.TenantID,
		AccountID:         account.AccountID,
		Name:              account.Name,
		Active:            account.Active,
		TagsCurrent:       account.TagsCurrent.Get(),
		TagsDesired:       account.TagsDesired.Get(),
		TagsDriftDetected: account.TagsDriftDetected,
		CreatedAt:         account.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         account.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewCloudAccountListResponse(accounts []db.CloudAccount) []render.Renderer {
	var list []render.Renderer
	for _, account := range accounts {
		list = append(list, NewCloudAccountResponse(account))
	}
	return list
}
