package middleware

import (
	"context"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

func APIKeyCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idString := chi.URLParam(r, "id")

			id, err := uuid.Parse(idString)
			if err != nil {
				_ = render.Render(w, r, responses.ErrNotFound)
				return
			}

			apiKey, err := app.Repository.FindAPIKey(r.Context(), id)
			if err != nil {
				if err == pgx.ErrNoRows {
					_ = render.Render(w, r, responses.ErrNotFound)
					return
				}

				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), ContextAPIKey, apiKey)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
