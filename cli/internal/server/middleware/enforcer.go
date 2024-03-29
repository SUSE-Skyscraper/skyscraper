package middleware

import (
	"fmt"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/cli/fga"

	"github.com/suse-skyscraper/skyscraper/cli/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
)

func EnforcerHandler(app *application.App, document fga.Document, relation fga.Relation) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			caller, ok := r.Context().Value(ContextCaller).(auth.Caller)
			if !ok {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			var objectID string
			switch document {
			case fga.DocumentAccount:
				// we need to convert group,tenantID,resourceID to ID
				group := chi.URLParam(r, "group")
				tenantID := chi.URLParam(r, "tenant_id")
				resourceID := chi.URLParam(r, "resource_id")
				resource, err := app.Repo.FindCloudAccountByCloudAndTenant(r.Context(), db.FindCloudAccountByCloudAndTenantParams{
					Cloud:     group,
					TenantID:  tenantID,
					AccountID: resourceID,
				})
				if err != nil {
					_ = render.Render(w, r, responses.ErrInternalServerError)
					return
				}

				objectID = resource.ID.String()
			case fga.DocumentOrganization:
				objectID = fga.DefaultOrganizationID
			}

			allowed, err := app.FGAClient.Check(r.Context(), caller.ID, relation, document, objectID)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			} else if !allowed {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
