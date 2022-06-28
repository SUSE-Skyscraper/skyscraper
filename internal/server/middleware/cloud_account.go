package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func CloudAccountCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := chi.URLParam(r, "tenant_id")
			cloudProvider := chi.URLParam(r, "cloud")
			id := chi.URLParam(r, "id")

			account, err := app.DB.GetCloudAccount(context.TODO(), db.GetCloudAccountParams{
				Cloud:     cloudProvider,
				TenantID:  tenantID,
				AccountID: id,
			})
			if err != nil {
				if err == pgx.ErrNoRows {
					_ = render.Render(w, r, responses.ErrNotFound)
					return
				}

				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), CloudAccount, account)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
