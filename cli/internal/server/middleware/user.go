package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
)

func UserCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")

			idParsed, err := uuid.Parse(id)
			if err != nil {
				_ = render.Render(w, r, responses.ErrNotFound)
				return
			}

			user, err := app.Repo.GetUser(r.Context(), idParsed)
			if err != nil {
				if err == pgx.ErrNoRows {
					_ = render.Render(w, r, responses.ErrNotFound)
					return
				}

				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUser, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
