package server

import (
	"context"
	"encoding/json"
	"fmt"
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
		filters := getCloudAccountFilters(r)

		for key, value := range r.URL.Query() {
			filters[key] = value[0]
		}

		cloudAccounts, err := app.Repository.SearchCloudAccounts(r.Context(), db.SearchCloudAccountsInput{
			Filters: filters,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudAccountListResponse(cloudAccounts))
	}
}

func V1UpdateCloudTenantAccount(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	natsWorker := workers.NewWorker(app)

	return func(w http.ResponseWriter, r *http.Request) {
		cloudAccount, ok := r.Context().Value(middleware.CloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		before, err := json.Marshal(cloudAccount)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		var payload payloads.UpdateCloudAccountPayload
		err = render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		defer func(repo db.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		account, err := repo.UpdateCloudAccount(r.Context(), db.UpdateCloudAccountParams{
			ID:          cloudAccount.ID,
			TagsDesired: payload.Data.GetJSON(),
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		after, err := json.Marshal(account)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_, err = repo.CreateAuditLog(r.Context(), db.CreateAuditLogParams{
			UserID:       user.ID,
			ResourceType: db.AuditResourceTypeCloudAccount,
			ResourceID:   cloudAccount.ID,
			Message:      fmt.Sprintf("Updated cloud account from %s to %s", string(before), string(after)),
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = natsWorker.PublishTagChange(workers.ChangeTagsPayload{
			ID:          account.ID.String(),
			AccountName: account.Name,
		})
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

func getCloudAccountFilters(r *http.Request) map[string]interface{} {
	cloudTenantID := chi.URLParam(r, "tenant_id")
	cloudProvider := chi.URLParam(r, "cloud")

	filters := make(map[string]interface{})
	if cloudTenantID != "" {
		filters["tenant_id"] = cloudTenantID
	}
	if cloudProvider != "" {
		filters["cloud"] = cloudProvider
	}

	for key, value := range r.URL.Query() {
		filters[key] = value[0]
	}

	return filters
}
