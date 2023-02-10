package payloads

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/pkg/errors"
)

//----------------------------------------------------------------------------------------------------------------------
// Assign Cloud Accounts to Organizational Units
//----------------------------------------------------------------------------------------------------------------------

// AssignCloudAccountToOUPayloadData is the data for the AssignCloudAccountToOUPayload.
type AssignCloudAccountToOUPayloadData struct {
	OrganizationalUnitID string `json:"organizational_unit_id"`
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

	// return an error if the uuid is not valid
	_, err := uuid.Parse(u.Data.OrganizationalUnitID)
	if err != nil {
		return err
	}

	return nil
}
