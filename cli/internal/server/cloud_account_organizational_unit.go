package server

import (
	"context"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/payloads"
	responses2 "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	db2 "github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/render"
)

func V1AssignCloudAccountToOU(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the cloud account that we're changing
		cloudAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db2.CloudAccount)
		if !ok {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.AssignCloudAccountToOUPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db2.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		err = repo.UnAssignCloudAccountFromOrganizationalUnits(r.Context(), cloudAccount.ID)
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		err = repo.AssignCloudAccountToOrganizationalUnit(r.Context(), cloudAccount.ID, payload.Data.GetOrganizationalUnitID())
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses2.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
