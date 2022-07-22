package responses

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type APIKeyItemAttributes struct {
	Token       string `json:"token,omitempty"`
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

func newAPIKeyItem(apiKey db.ApiKey, token string) APIKeyItem {
	return APIKeyItem{
		ID:   apiKey.ID.String(),
		Type: ObjectResponseTypeAPIKey,
		Attributes: APIKeyItemAttributes{
			Owner:       apiKey.Owner,
			Description: apiKey.Description.String,
			System:      apiKey.System,
			Token:       token,
			CreatedAt:   apiKey.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   apiKey.UpdatedAt.Format(time.RFC3339),
		},
	}
}

type APIKeyResponse struct {
	Data APIKeyItem `json:"data"`
}

func (rd *APIKeyResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewAPIKeyResponse(apiKey db.ApiKey, token string) *APIKeyResponse {
	return &APIKeyResponse{
		Data: newAPIKeyItem(apiKey, token),
	}
}

type APIKeysResponse struct {
	Data []APIKeyItem `json:"data"`
}

func (rd *APIKeysResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewAPIKeysResponse(apiKeys []db.ApiKey) *APIKeysResponse {
	list := make([]APIKeyItem, len(apiKeys))
	for i, key := range apiKeys {
		list[i] = newAPIKeyItem(key, "")
	}

	return &APIKeysResponse{
		Data: list,
	}
}
