package responses

import (
	"time"

	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

func newAPIKeyItem(apiKey db.ApiKey, token string) resp.APIKeyItem {
	return resp.APIKeyItem{
		ID:   apiKey.ID.String(),
		Type: resp.ObjectResponseTypeAPIKey,
		Attributes: resp.APIKeyItemAttributes{
			Owner:       apiKey.Owner,
			Description: apiKey.Description.String,
			System:      apiKey.System,
			Token:       token,
			CreatedAt:   apiKey.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   apiKey.UpdatedAt.Format(time.RFC3339),
		},
	}
}

func NewAPIKeyResponse(apiKey db.ApiKey, token string) *resp.APIKeyResponse {
	return &resp.APIKeyResponse{
		Data: newAPIKeyItem(apiKey, token),
	}
}

func NewAPIKeysResponse(apiKeys []db.ApiKey) *resp.APIKeysResponse {
	list := make([]resp.APIKeyItem, len(apiKeys))
	for i, key := range apiKeys {
		list[i] = newAPIKeyItem(key, "")
	}

	return &resp.APIKeysResponse{
		Data: list,
	}
}
