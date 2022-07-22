package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1Profile(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := r.Context().Value(middleware.ContextCaller).(auth.Caller)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Only show the profile for users
		if caller.Type != auth.CallerUser {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		user, err := app.Repository.FindUser(r.Context(), caller.ID.String())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewUserResponse(user))
	}
}
