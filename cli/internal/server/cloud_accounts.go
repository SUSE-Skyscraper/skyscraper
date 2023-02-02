package server

import (
	"context"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/queue"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/api/payloads"
	responses2 "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/responses"
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
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudAccountListResponse(cloudAccounts))
	}
}

func V1UpdateCloudAccount(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	natsWorker := queue.NewPluginWorker(app)

	return func(w http.ResponseWriter, r *http.Request) {
		// Get the cloud account that we're changing
		cloudAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.UpdateCloudAccountPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
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
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// AuditChange the change
		err = auditor.AuditChange(r.Context(), db.AuditResourceTypeCloudAccount, account.ID, payload)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// Publish the change to the NATS queue.
		// If this fails, we don't care because it can be retried later.
		// It's more important that we update the account.
		_ = natsWorker.PublishMessage(account.Cloud, queue.PluginPayload{
			ResourceID: account.AccountID,
			Cloud:      account.Cloud,
			TenantID:   account.TenantID,
			Action:     queue.PluginActionTagUpdate,
		})

		_ = render.Render(w, r, responses.NewCloudAccountResponse(account))
	}
}

func V1GetCloudAccount(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
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
