package middleware

import (
	"fmt"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/fga"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
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
				objectID = chi.URLParam(r, "id") // will deprecate soon
				if objectID == "" {
					objectID = chi.URLParam(r, "resource_id")
				}
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
