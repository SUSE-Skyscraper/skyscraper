package payloads

import (
	"net/http"

	"github.com/google/uuid"
)

type CreateOrganizationalUnitsPayloadData struct {
	ParentID    string `json:"parent_id"`
	DisplayName string `json:"display_name"`
}

type CreateOrganizationalUnitsPayload struct {
	Data CreateOrganizationalUnitsPayloadData `json:"data"`
}

func (u *CreateOrganizationalUnitsPayload) Bind(_ *http.Request) error {
	// we can allow an empty parent id, but when we have a value, it should be valid
	if u.Data.ParentID != "" {
		_, err := uuid.Parse(u.Data.ParentID)
		if err != nil {
			return err
		}
	}

	return nil
}
