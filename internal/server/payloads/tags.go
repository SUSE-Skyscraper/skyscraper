package payloads

import (
	"net/http"
)

type UpdateTagPayloadData struct {
	DisplayName string `json:"display_name"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

type UpdateTagPayload struct {
	Data UpdateTagPayloadData `json:"data"`
}

func (u *UpdateTagPayload) Bind(_ *http.Request) error {
	return nil
}

type CreateTagPayloadData struct {
	DisplayName string `json:"display_name"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
	Key         string `json:"key"`
}

type CreateTagPayload struct {
	Data CreateTagPayloadData `json:"data"`
}

func (u *CreateTagPayload) Bind(_ *http.Request) error {
	return nil
}
