package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/suse-skyscraper/skyscraper/api/responses"

	"github.com/jackc/pgx/v4"

	"github.com/suse-skyscraper/skyscraper/api/queue"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgtype"
	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
)

func V1GetResource(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewCloudAccountResponse(cloudTenantAccount))
	}
}

func V1ListResources(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseAccountSearchFilters(r)

		for key, value := range r.URL.Query() {
			filters[key] = value[0]
		}

		group := chi.URLParam(r, "group")
		if group == "" {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}
		filters["cloud"] = group

		tenantID := chi.URLParam(r, "tenant_id")
		if tenantID == "" {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}
		filters["tenant_id"] = tenantID

		cloudAccounts, err := app.Searcher.SearchCloudAccounts(r.Context(), db.SearchCloudAccountsInput{
			Filters: filters,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewCloudAccountListResponse(cloudAccounts))
	}
}

func V1CreateOrUpdateResource(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	natsWorker := queue.NewPluginWorker(app)

	return func(w http.ResponseWriter, r *http.Request) {
		resourceID := chi.URLParam(r, "resource_id")
		if resourceID == "" {
			_ = render.Render(w, r, responses.ErrInvalidRequest(fmt.Errorf("resource_id is required")))
			return
		}

		tenant, ok := r.Context().Value(middleware.ContextTenant).(db.CloudTenant)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.CreateOrUpdateResourcePayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		tx, err := app.PostgresPool.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}
		repo := app.Repo.WithTx(tx)

		// If any error occurs, rollback the transaction
		// We'll call this even if the transaction is committed, because it's a no-op if the transaction is already committed.
		defer func(tx pgx.Tx, ctx context.Context) {
			_ = tx.Rollback(ctx)
		}(tx, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		// Set the jsonb field. As we're certain of the type, we can be sure it will succeed.
		tagsDesired := pgtype.JSONB{}
		_ = tagsDesired.Set(payload.Data.TagsDesired)

		// Set the jsonb field. As we're certain of the type, we can be sure it will succeed.
		tagsCurrent := pgtype.JSONB{}
		_ = tagsCurrent.Set(payload.Data.TagsCurrent)

		account, err := repo.CreateOrUpdateCloudAccount(r.Context(), db.CreateOrUpdateCloudAccountParams{
			TenantID:    tenant.TenantID,
			Cloud:       tenant.Cloud,
			AccountID:   resourceID,
			Name:        payload.Data.AccountName,
			TagsCurrent: tagsCurrent,
			TagsDesired: tagsDesired,
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

		err = app.FGAClient.AddAccountToOrganization(r.Context(), account.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = tx.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// if we change the body desired, and it's not the same as the body current, then publish the change to nats
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

		_ = render.Render(w, r, NewCloudAccountResponse(account))
	}
}

func parseAccountSearchFilters(r *http.Request) map[string]interface{} {
	filters := make(map[string]interface{})

	for key, value := range r.URL.Query() {
		filters[key] = value[0]
	}

	return filters
}

// newCloudAccount creates a new CloudAccountItem from a db.CloudAccount.
func newCloudAccount(account db.CloudAccount) responses.CloudAccountItem {
	var currentTags map[string]string
	_ = json.Unmarshal(account.TagsCurrent.Bytes, &currentTags)
	var desiredTags map[string]string
	_ = json.Unmarshal(account.TagsDesired.Bytes, &desiredTags)

	return responses.CloudAccountItem{
		ID:   account.ID.String(),
		Type: responses.ObjectResponseTypeCloudAccount,
		Attributes: responses.CloudAccountItemAttributes{
			CloudProvider:     account.Cloud,
			TenantID:          account.TenantID,
			AccountID:         account.AccountID,
			Name:              account.Name,
			Active:            account.Active,
			TagsCurrent:       currentTags,
			TagsDesired:       desiredTags,
			TagsDriftDetected: account.TagsDriftDetected,
			CreatedAt:         account.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         account.UpdatedAt.Format(time.RFC3339),
		},
	}
}

// NewCloudAccountResponse creates a new resp.CloudAccountResponse from a db.CloudAccount.
func NewCloudAccountResponse(account db.CloudAccount) *responses.CloudAccountResponse {
	cloudAccount := newCloudAccount(account)
	return &responses.CloudAccountResponse{
		Data: cloudAccount,
	}
}

// NewCloudAccountListResponse creates a new resp.CloudAccountListResponse from a list of db.CloudAccount.
func NewCloudAccountListResponse(accounts []db.CloudAccount) *responses.CloudAccountListResponse {
	list := make([]responses.CloudAccountItem, len(accounts))
	for i, account := range accounts {
		list[i] = newCloudAccount(account)
	}

	return &responses.CloudAccountListResponse{
		Data: list,
	}
}
