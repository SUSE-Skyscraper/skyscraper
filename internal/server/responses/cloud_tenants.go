package responses

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type CloudTenantResponse struct {
	CloudProvider string `json:"cloud_provider"`
	TenantID      string `json:"tenant_id"`
	Name          string `json:"name"`
	Active        bool   `json:"active"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

func (rd *CloudTenantResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
func NewCloudTenantResponse(tenant db.CloudTenant) *CloudTenantResponse {
	return &CloudTenantResponse{
		CloudProvider: tenant.Cloud,
		TenantID:      tenant.TenantID,
		Name:          tenant.Name,
		Active:        tenant.Active,
		CreatedAt:     tenant.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     tenant.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewCloudTenantListResponse(tenants []db.CloudTenant) []render.Renderer {
	var list []render.Renderer
	for _, tenant := range tenants {
		list = append(list, NewCloudTenantResponse(tenant))
	}
	return list
}
