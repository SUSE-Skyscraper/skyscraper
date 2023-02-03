package responses

import "net/http"

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

type APIKeyResponse struct {
	Data APIKeyItem `json:"data"`
}

func (rd *APIKeyResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type APIKeysResponse struct {
	Data []APIKeyItem `json:"data"`
}

func (rd *APIKeysResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
