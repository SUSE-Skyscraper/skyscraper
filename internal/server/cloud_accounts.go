package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
	"github.com/suse-skyscraper/skyscraper/workers"
)

func V1ListCloudAccounts(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantID := chi.URLParam(r, "tenant_id")
		cloudProvider := chi.URLParam(r, "cloud")

		var cloudTenantAccounts, err = app.DB.GetCloudAllAccountsForCloudAndTenant(
			r.Context(),
			db.GetCloudAllAccountsForCloudAndTenantParams{
				Cloud:    cloudProvider,
				TenantID: cloudTenantID,
			})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.RenderList(w, r, responses.NewCloudAccountListResponse(cloudTenantAccounts))
	}
}

func V1UpdateCloudTenantAccount(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	natsWorker := workers.NewWorker(app)

	return func(w http.ResponseWriter, r *http.Request) {
		tenantID := chi.URLParam(r, "tenant_id")
		cloudProvider := chi.URLParam(r, "cloud")
		id := chi.URLParam(r, "id")

		var payload payloads.UpdateCloudAccountPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		err = app.DB.UpdateCloudAccount(r.Context(), db.UpdateCloudAccountParams{
			Cloud:       cloudProvider,
			TenantID:    tenantID,
			AccountID:   id,
			TagsDesired: payload.GetJSON(),
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		account, err := app.DB.GetCloudAccount(r.Context(), db.GetCloudAccountParams{
			Cloud:     cloudProvider,
			TenantID:  tenantID,
			AccountID: id,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		changeCloudPayload := workers.ChangeTagsPayload{
			Cloud:       cloudProvider,
			TenantID:    tenantID,
			AccountID:   id,
			AccountName: account.Name,
		}
		err = natsWorker.PublishTagChange(changeCloudPayload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudAccountResponse(account))
	}
}

func V1GetCloudAccount(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantAccount := r.Context().Value(middleware.CloudAccount).(db.CloudAccount)

		_ = render.Render(w, r, responses.NewCloudAccountResponse(cloudTenantAccount))
	}
}
