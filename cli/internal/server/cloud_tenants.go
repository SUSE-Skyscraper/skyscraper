package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/api/responses"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/db"
)

func V1ListCloudTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenants, err := app.Repo.GetCloudTenants(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewCloudTenantListResponse(cloudTenants))
	}
}

func V1CreateOrUpdateTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chi.URLParam(r, "tenant_id")
		if tenantID == "" {
			_ = render.Render(w, r, responses.ErrInvalidRequest(fmt.Errorf("id is required")))
			return
		}

		// Bind the payload
		var payload payloads.CreateOrUpdateTenantPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		tenant, err := app.Repo.CreateOrUpdateCloudTenant(r.Context(), db.CreateOrUpdateCloudTenantParams{
			Cloud:    payload.Data.Cloud,
			TenantID: tenantID,
			Name:     payload.Data.Name,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewCloudTenantResponse(tenant))
	}
}

//----------------------------------------------------------------------------------------------------------------------
// Common
//----------------------------------------------------------------------------------------------------------------------

// newCloudTenantItem creates a new CloudTenantItem from a db.CloudTenant.
func newCloudTenantItem(tenant db.CloudTenant) responses.CloudTenantItem {
	return responses.CloudTenantItem{
		ID:   tenant.ID.String(),
		Type: responses.ObjectResponseTypeCloudTenant,
		Attributes: responses.CloudTenantAttributes{
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
func NewCloudTenantResponse(tenant db.CloudTenant) *responses.CloudTenantResponse {
	data := newCloudTenantItem(tenant)

	return &responses.CloudTenantResponse{
		Data: data,
	}
}

//----------------------------------------------------------------------------------------------------------------------
// List
//----------------------------------------------------------------------------------------------------------------------

func NewCloudTenantListResponse(tenants []db.CloudTenant) *responses.CloudTenantsResponse {
	list := make([]responses.CloudTenantItem, len(tenants))
	for i, tenant := range tenants {
		list[i] = newCloudTenantItem(tenant)
	}
	return &responses.CloudTenantsResponse{Data: list}
}
