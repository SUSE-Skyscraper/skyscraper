package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/suse-skyscraper/skyscraper/api/payloads"
	responses2 "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/responses"

	"github.com/go-chi/render"
)

func V1ListCloudTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenants, err := app.Repository.GetCloudTenants(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudTenantListResponse(cloudTenants))
	}
}

func V1CreateOrUpdateTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chi.URLParam(r, "tenant_id")
		if tenantID == "" {
			_ = render.Render(w, r, responses2.ErrInvalidRequest(fmt.Errorf("id is required")))
			return
		}

		// Bind the payload
		var payload payloads.CreateOrUpdateTenantPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInvalidRequest(err))
			return
		}

		tenant, err := app.Repository.CreateOrUpdateCloudTenant(r.Context(), db.CreateOrUpdateCloudTenantParams{
			Cloud:    payload.Data.Cloud,
			TenantID: tenantID,
			Name:     payload.Data.Name,
		})
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudTenantResponse(tenant))
	}
}
