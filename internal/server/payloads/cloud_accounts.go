package payloads

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type UpdateCloudAccountPayloadData struct {
	TagsDesired map[string]string `json:"tags_desired"`
	json        pgtype.JSONB
}

type UpdateCloudAccountPayload struct {
	Data UpdateCloudAccountPayloadData `json:"data"`
}

func (u *UpdateCloudAccountPayload) Bind(_ *http.Request) error {
	if u.Data.TagsDesired == nil {
		return errors.Errorf("tags_desired is required")
	}

	jsonTags := pgtype.JSONB{}
	err := jsonTags.Set(u.Data.TagsDesired)
	if err != nil {
		return err
	}
	u.Data.json = jsonTags

	return nil
}

func (u *UpdateCloudAccountPayloadData) GetJSON() pgtype.JSONB {
	return u.json
}

type AssignCloudAccountToOUPayloadData struct {
	OrganizationalUnitID string `json:"organizational_unit_id"`
	organizationalUnitID uuid.UUID
}

type AssignCloudAccountToOUPayload struct {
	Data AssignCloudAccountToOUPayloadData `json:"data"`
}

func (u *AssignCloudAccountToOUPayload) Bind(_ *http.Request) error {
	id, err := uuid.Parse(u.Data.OrganizationalUnitID)
	if err != nil {
		return err
	}
	u.Data.organizationalUnitID = id

	return nil
}

func (d *AssignCloudAccountToOUPayloadData) GetOrganizationalUnitID() uuid.UUID {
	return d.organizationalUnitID
}
