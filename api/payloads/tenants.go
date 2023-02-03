package payloads

import (
	"fmt"
	"net/http"
)

//----------------------------------------------------------------------------------------------------------------------
// Create or Update Tenants
//----------------------------------------------------------------------------------------------------------------------

// CreateOrUpdateTenantPayloadData is the data for the CreateOrUpdateTenantPayload.
type CreateOrUpdateTenantPayloadData struct {
	Name  string `json:"name"`
	Cloud string `json:"cloud"`
}

// CreateOrUpdateTenantPayload is the payload for creating a tenant.
type CreateOrUpdateTenantPayload struct {
	Data CreateOrUpdateTenantPayloadData `json:"data"`
}

// Bind binds extra data from the payload CreateOrUpdateResourcePayload.
func (u *CreateOrUpdateTenantPayload) Bind(_ *http.Request) error {
	if u.Data.Name == "" {
		return fmt.Errorf("name is required")
	} else if u.Data.Cloud == "" {
		return fmt.Errorf("cloud is required")
	}

	return nil
}
