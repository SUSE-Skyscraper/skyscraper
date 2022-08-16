package payloads

import (
	"net/http"

	"github.com/google/uuid"
)

type CreateOrganizationalUnitsPayloadData struct {
	ParentID       string `json:"parent_id"`
	DisplayName    string `json:"display_name"`
	parentIDParsed uuid.NullUUID
}

type CreateOrganizationalUnitsPayload struct {
	Data CreateOrganizationalUnitsPayloadData `json:"data"`
}

func (u *CreateOrganizationalUnitsPayload) Bind(_ *http.Request) error {
	if u.Data.ParentID != "" {
		id, err := uuid.Parse(u.Data.ParentID)
		if err != nil {
			return err
		}
		u.Data.parentIDParsed = uuid.NullUUID{UUID: id, Valid: true}
	} else {
		u.Data.parentIDParsed = uuid.NullUUID{Valid: false}
	}

	return nil
}

func (d *CreateOrganizationalUnitsPayloadData) GetParentID() uuid.NullUUID {
	return d.parentIDParsed
}
