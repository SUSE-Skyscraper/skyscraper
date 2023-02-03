package payloads

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

//----------------------------------------------------------------------------------------------------------------------
// Update Cloud Accounts
//----------------------------------------------------------------------------------------------------------------------

// UpdateCloudAccountPayloadData is the data for the UpdateCloudAccountPayload.
type UpdateCloudAccountPayloadData struct {
	TagsDesired map[string]string `json:"tags_desired"`
	json        pgtype.JSONB
}

// UpdateCloudAccountPayload is the payload for updating a cloud account.
type UpdateCloudAccountPayload struct {
	Data UpdateCloudAccountPayloadData `json:"data"`
}

// Bind binds extra data from the payload UpdateCloudAccountPayload.
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

// GetJSON returns the parsed JSONB for the tags.
func (u *UpdateCloudAccountPayloadData) GetJSON() pgtype.JSONB {
	return u.json
}

//----------------------------------------------------------------------------------------------------------------------
// Assign Cloud Accounts to Organizational Units
//----------------------------------------------------------------------------------------------------------------------

// AssignCloudAccountToOUPayloadData is the data for the AssignCloudAccountToOUPayload.
type AssignCloudAccountToOUPayloadData struct {
	OrganizationalUnitID string `json:"organizational_unit_id"`
	organizationalUnitID uuid.UUID
}

// AssignCloudAccountToOUPayload is the payload for assigning a cloud account to an organizational unit.
type AssignCloudAccountToOUPayload struct {
	Data AssignCloudAccountToOUPayloadData `json:"data"`
}

// Bind binds extra data from the payload AssignCloudAccountToOUPayload.
func (u *AssignCloudAccountToOUPayload) Bind(_ *http.Request) error {
	if u.Data.OrganizationalUnitID == "" {
		return errors.Errorf("organizational_unit_id is required")
	}

	id, err := uuid.Parse(u.Data.OrganizationalUnitID)
	if err != nil {
		return err
	}
	u.Data.organizationalUnitID = id

	return nil
}

// GetOrganizationalUnitID returns the organizational unit ID.
func (d *AssignCloudAccountToOUPayloadData) GetOrganizationalUnitID() uuid.UUID {
	return d.organizationalUnitID
}
