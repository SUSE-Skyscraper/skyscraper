package responses

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type CloudTenantAttributes struct {
	TenantID      string `json:"tenant_id"`
	CloudProvider string `json:"cloud_provider"`
	Name          string `json:"name"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type CloudTenantItem struct {
	ID         string                `json:"id"`
	Type       ObjectResponseType    `json:"type"`
	Attributes CloudTenantAttributes `json:"attributes"`
}

type CloudTenantsResponse struct {
	Data []CloudTenantItem `json:"data"`
}

func (rd *CloudTenantsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewCloudTenantListResponse(tenants []db.CloudTenant) *CloudTenantsResponse {
	list := make([]CloudTenantItem, len(tenants))
	for i, tenant := range tenants {
		list[i] = newCloudTenantItem(tenant)
	}
	return &CloudTenantsResponse{Data: list}
}

func newCloudTenantItem(tenant db.CloudTenant) CloudTenantItem {
	return CloudTenantItem{
		ID:   tenant.ID.String(),
		Type: ObjectResponseTypeCloudTenant,
		Attributes: CloudTenantAttributes{
			CloudProvider: tenant.Cloud,
			TenantID:      tenant.TenantID,
			Name:          tenant.Name,
			Active:        tenant.Active,
			CreatedAt:     tenant.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     tenant.UpdatedAt.Format(time.RFC3339),
		},
	}
}
