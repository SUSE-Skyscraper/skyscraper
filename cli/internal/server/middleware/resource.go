package middleware

import (
	"context"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
)

func ResourceCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := chi.URLParam(r, "tenant_id")
			group := chi.URLParam(r, "group")
			resourceID := chi.URLParam(r, "resource_id")

			cloudAccount, err := app.Repository.FindCloudAccount(r.Context(), db.FindCloudAccountInput{
				Cloud:     group,
				TenantID:  tenantID,
				AccountID: resourceID,
			})
			if err != nil {
				if err == pgx.ErrNoRows {
					_ = render.Render(w, r, responses.ErrNotFound)
					return
				}

				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), ContextCloudAccount, cloudAccount)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
