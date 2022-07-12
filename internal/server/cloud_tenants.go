package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1CloudTenants(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cloudTenants, err := app.DB.GetCloudTenants(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewCloudTenantListResponse(cloudTenants))
	}
}

func V1CloudTenantTags(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags := []string{"Owner", "Group", "PoNumber", "CostCenter", "Stakeholder", "Department", "GeneralLedgerCode",
			"Environment", "FinanceBusinessPartner", "BillingContacts", "AdminContacts", "SlackChannels"}

		_ = render.Render(w, r, responses.NewCloudAccountTagsResponse(tags))
	}
}
