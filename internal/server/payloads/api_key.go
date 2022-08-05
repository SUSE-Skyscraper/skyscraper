package payloads

import (
	"net/http"
)

type CreateAPIKeyPayloadData struct {
	Owner       string `json:"owner"`
	Description string `json:"description"`
}

type CreateAPIKeyPayload struct {
	Data CreateAPIKeyPayloadData `json:"data"`
}

func (u *CreateAPIKeyPayload) Bind(_ *http.Request) error {
	return nil
}
