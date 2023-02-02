package responses

import (
	"time"

	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

//----------------------------------------------------------------------------------------------------------------------
// Common
//----------------------------------------------------------------------------------------------------------------------

// newCloudTenantItem creates a new CloudTenantItem from a db.CloudTenant.
func newCloudTenantItem(tenant db.CloudTenant) resp.CloudTenantItem {
	return resp.CloudTenantItem{
		ID:   tenant.ID.String(),
		Type: resp.ObjectResponseTypeCloudTenant,
		Attributes: resp.CloudTenantAttributes{
			CloudProvider: tenant.Cloud,
			TenantID:      tenant.TenantID,
			Name:          tenant.Name,
			Active:        tenant.Active,
			CreatedAt:     tenant.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     tenant.UpdatedAt.Format(time.RFC3339),
		},
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Single
//----------------------------------------------------------------------------------------------------------------------

// NewCloudTenantResponse creates a new CloudTenantResponse from a db.CloudTenant.
func NewCloudTenantResponse(tenant db.CloudTenant) *resp.CloudTenantResponse {
	data := newCloudTenantItem(tenant)

	return &resp.CloudTenantResponse{
		Data: data,
	}
}

//----------------------------------------------------------------------------------------------------------------------
// List
//----------------------------------------------------------------------------------------------------------------------

func NewCloudTenantListResponse(tenants []db.CloudTenant) *resp.CloudTenantsResponse {
	list := make([]resp.CloudTenantItem, len(tenants))
	for i, tenant := range tenants {
		list[i] = newCloudTenantItem(tenant)
	}
	return &resp.CloudTenantsResponse{Data: list}
}
