package responses

import (
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type APIKeyItemAttributes struct {
	Owner       string `json:"owner"`
	Description string `json:"description"`
	System      bool   `json:"system"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type APIKeyItem struct {
	ID         string               `json:"id"`
	Type       ObjectResponseType   `json:"type"`
	Attributes APIKeyItemAttributes `json:"attributes"`
}

func newAPIKeyItem(apiKey db.ApiKey) APIKeyItem {
	return APIKeyItem{
		ID:   apiKey.ID.String(),
		Type: ObjectResponseTypeAPIKey,
		Attributes: APIKeyItemAttributes{
			Owner:       apiKey.Owner,
			Description: apiKey.Description.String,
			System:      apiKey.System,
			CreatedAt:   apiKey.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   apiKey.UpdatedAt.Format(time.RFC3339),
		},
	}
}
