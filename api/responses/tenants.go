package responses

import (
	"net/http"
)

//----------------------------------------------------------------------------------------------------------------------
// Common
//----------------------------------------------------------------------------------------------------------------------

// CloudTenantAttributes represents the attributes of a cloud tenant.
type CloudTenantAttributes struct {
	TenantID      string `json:"tenant_id"`
	CloudProvider string `json:"cloud_provider"`
	Name          string `json:"name"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CloudTenantItem represents a cloud tenant.
type CloudTenantItem struct {
	ID         string                `json:"id"`
	Type       ObjectResponseType    `json:"type"`
	Attributes CloudTenantAttributes `json:"attributes"`
}

//----------------------------------------------------------------------------------------------------------------------
// Single
//----------------------------------------------------------------------------------------------------------------------

// CloudTenantResponse represents a single tenant.
type CloudTenantResponse struct {
	Data CloudTenantItem `json:"data"`
}

// Render is a no-op.
func (rd *CloudTenantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// List
//----------------------------------------------------------------------------------------------------------------------

// CloudTenantsResponse represents a list of cloud tenants.
type CloudTenantsResponse struct {
	Data []CloudTenantItem `json:"data"`
}

// Render is a noop.
func (rd *CloudTenantsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
