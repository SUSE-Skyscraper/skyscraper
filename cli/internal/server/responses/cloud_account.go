package responses

import (
	"encoding/json"
	"time"

	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

// newCloudAccount creates a new CloudAccountItem from a db.CloudAccount.
func newCloudAccount(account db.CloudAccount) resp.CloudAccountItem {
	var currentTags map[string]string
	_ = json.Unmarshal(account.TagsCurrent.Bytes, &currentTags)
	var desiredTags map[string]string
	_ = json.Unmarshal(account.TagsDesired.Bytes, &desiredTags)

	return resp.CloudAccountItem{
		ID:   account.ID.String(),
		Type: resp.ObjectResponseTypeCloudAccount,
		Attributes: resp.CloudAccountItemAttributes{
			CloudProvider:     account.Cloud,
			TenantID:          account.TenantID,
			AccountID:         account.AccountID,
			Name:              account.Name,
			Active:            account.Active,
			TagsCurrent:       currentTags,
			TagsDesired:       desiredTags,
			TagsDriftDetected: account.TagsDriftDetected,
			CreatedAt:         account.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         account.UpdatedAt.Format(time.RFC3339),
		},
	}
}

// NewCloudAccountResponse creates a new resp.CloudAccountResponse from a db.CloudAccount.
func NewCloudAccountResponse(account db.CloudAccount) *resp.CloudAccountResponse {
	cloudAccount := newCloudAccount(account)
	return &resp.CloudAccountResponse{
		Data: cloudAccount,
	}
}

// NewCloudAccountListResponse creates a new resp.CloudAccountListResponse from a list of db.CloudAccount.
func NewCloudAccountListResponse(accounts []db.CloudAccount) *resp.CloudAccountListResponse {
	list := make([]resp.CloudAccountItem, len(accounts))
	for i, account := range accounts {
		list[i] = newCloudAccount(account)
	}

	return &resp.CloudAccountListResponse{
		Data: list,
	}
}
