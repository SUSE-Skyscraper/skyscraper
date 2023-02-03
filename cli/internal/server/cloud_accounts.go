package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/api/responses"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
)

// todo: move to resource
func V1ListCloudAccounts(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseAccountSearchFilters(r)

		for key, value := range r.URL.Query() {
			filters[key] = value[0]
		}

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

// todo: move to resource
func V1GetCloudAccount(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenantAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewCloudAccountResponse(cloudTenantAccount))
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
