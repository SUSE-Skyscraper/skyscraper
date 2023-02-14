package middleware

import (
	"context"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
)

func TenantCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			group := chi.URLParam(r, "group")
			tenantID := chi.URLParam(r, "tenant_id")

			tenant, err := app.Repo.GetCloudTenant(r.Context(), db.GetCloudTenantParams{
				Cloud:    group,
				TenantID: tenantID,
			})
			if err != nil {
				if err == pgx.ErrNoRows {
					_ = render.Render(w, r, responses.ErrNotFound)
					return
				}

				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), ContextTenant, tenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
