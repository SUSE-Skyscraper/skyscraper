package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/auth"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func V1CallerProfile(app *application.App) func(w http.ResponseWriter, r *http.Request) {
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

		user, err := app.Repo.GetUser(r.Context(), caller.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewUserResponse(user))
	}
}

func V1CallerCloudAccounts(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		caller, ok := r.Context().Value(middleware.ContextCaller).(auth.Caller)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		organizationalUnits, err := callerOrganizationalUnits(r.Context(), app, caller)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		ids := make([]uuid.UUID, 0, len(organizationalUnits))
		for _, ou := range organizationalUnits {
			ids = append(ids, ou.ID)
		}

		cloudAccounts, err := app.Repo.OrganizationalUnitsCloudAccounts(r.Context(), ids)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewCloudAccountListResponse(cloudAccounts))
	}
}

func callerOrganizationalUnits(ctx context.Context, app *application.App, caller auth.Caller) ([]db.OrganizationalUnit, error) {
	if caller.Type == auth.CallerUser {
		return app.Repo.GetUserOrganizationalUnits(ctx, caller.ID)
	} else if caller.Type == auth.CallerAPIKey {
		return app.Repo.GetAPIKeysOrganizationalUnits(ctx, caller.ID)
	}

	return nil, fmt.Errorf("caller not recognized")
}
