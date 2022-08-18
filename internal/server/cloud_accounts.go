package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
	"github.com/suse-skyscraper/skyscraper/workers"
)

func V1ListCloudAccounts(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseAccountSearchFilters(r)

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

func V1AssignCloudAccountToOU(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the cloud account that we're changing
		cloudAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.AssignCloudAccountToOUPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		err = repo.UnAssignCloudAccountFromOrganizationalUnits(r.Context(), cloudAccount.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = repo.AssignCloudAccountToOrganizationalUnit(r.Context(), cloudAccount.ID, payload.Data.GetOrganizationalUnitID())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func V1UpdateCloudAccount(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	natsWorker := workers.NewWorker(app)

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the cloud account that we're changing
		cloudAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.UpdateCloudAccountPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		// Update the cloud account
		account, err := repo.UpdateCloudAccount(r.Context(), db.UpdateCloudAccountParams{
			ID:          cloudAccount.ID,
			TagsDesired: payload.Data.GetJSON(),
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// AuditChange the change
		err = auditor.AuditChange(r.Context(), db.AuditResourceTypeCloudAccount, account.ID, payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Publish the change to the NATS queue.
		// If this fails, we don't care because it can be retried later.
		// It's more important that we update the account.
		_ = natsWorker.PublishTagChange(workers.ChangeTagsPayload{
			ID:          account.ID.String(),
			AccountName: account.Name,
		})

		_ = render.Render(w, r, responses.NewCloudAccountResponse(account))
	}
}

func V1GetCloudAccount(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudAccountResponse(cloudTenantAccount))
	}
}

func parseAccountSearchFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	for key, value := range r.URL.Query() {
		filters[key] = value[0]
	}

	return filters
}
