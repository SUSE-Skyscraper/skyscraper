package server

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/suse-skyscraper/skyscraper/api/queue"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgtype"
	"github.com/suse-skyscraper/skyscraper/api/payloads"
	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/responses"
)

func V1CreateOrUpdateResource(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	natsWorker := queue.NewPluginWorker(app)

	return func(w http.ResponseWriter, r *http.Request) {
		resourceID := chi.URLParam(r, "resource_id")
		if resourceID == "" {
			_ = render.Render(w, r, resp.ErrInvalidRequest(fmt.Errorf("resource_id is required")))
			return
		}

		tenant, ok := r.Context().Value(middleware.ContextTenant).(db.CloudTenant)
		if !ok {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.CreateOrUpdateResourcePayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, resp.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		account, err := app.Repository.CreateOrUpdateCloudAccount(r.Context(), db.CreateOrUpdateCloudAccountParams{
			TenantID:    tenant.TenantID,
			Cloud:       tenant.Cloud,
			AccountID:   resourceID,
			Name:        payload.Data.AccountName,
			TagsCurrent: payload.Data.GetTagsCurrent(),
			TagsDesired: payload.Data.GetTagsDesired(),
		})
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// AuditChange the change
		err = auditor.AuditChange(r.Context(), db.AuditResourceTypeCloudAccount, account.ID, payload)
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		err = app.FGAClient.AddAccountToOrganization(r.Context(), account.ID)
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// if we change the tags desired, and it's not the same as the tags current, then publish the change to nats
		if account.TagsDesired.Status == pgtype.Present && !reflect.DeepEqual(account.TagsDesired, account.TagsCurrent) {
			// Publish the change to the NATS queue.
			// If this fails, we don't care because it can be retried later.
			// It's more important that we update the account.
			_ = natsWorker.PublishMessage(account.Cloud, queue.PluginPayload{
				ResourceID: account.AccountID,
				Cloud:      account.Cloud,
				TenantID:   account.TenantID,
				Action:     queue.PluginActionTagUpdate,
			})
		}

		_ = render.Render(w, r, responses.NewCloudAccountResponse(account))
	}
}
