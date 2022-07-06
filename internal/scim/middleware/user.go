package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/scim/responses"
)

func UserCtx(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idString := chi.URLParam(r, "id")
			id, err := strconv.ParseInt(idString, 10, 32)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInvalidRequest(err))
				return
			}

			user, err := app.DB.GetUser(r.Context(), int32(id))
			if errors.Is(err, pgx.ErrNoRows) {
				_ = render.Render(w, r, responses.ErrNotFound)
				return
			} else if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), User, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
