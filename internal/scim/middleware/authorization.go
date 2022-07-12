package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func BearerAuthorizationHandler(app *application.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(authHeader) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized: authorization header malformed")
				return
			}

			token := authHeader[1]
			_, err := app.Repository.FindAPIKey(r.Context(), token)
			if err != nil && err == pgx.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprintf(w, "Not Authorized")
				return
			} else if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
