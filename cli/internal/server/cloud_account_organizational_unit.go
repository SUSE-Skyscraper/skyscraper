package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4"

	"github.com/suse-skyscraper/skyscraper/api/payloads"
	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/render"
)

func V1AssignCloudAccountToOU(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the cloud account that we're changing
		cloudAccount, ok := r.Context().Value(middleware.ContextCloudAccount).(db.CloudAccount)
		if !ok {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// Bind the payload
		var payload payloads.AssignCloudAccountToOUPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, resp.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		tx, err := app.PostgresPool.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}
		repo := app.Repo.WithTx(tx)

		// If any error occurs, rollback the transaction
		defer func(tx pgx.Tx, ctx context.Context) {
			_ = tx.Rollback(ctx)
		}(tx, r.Context())

		err = repo.UnAssignAccountFromOUs(r.Context(), cloudAccount.ID)
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		ouID, err := uuid.Parse(payload.Data.OrganizationalUnitID)
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		err = repo.AssignAccountToOU(r.Context(), db.AssignAccountToOUParams{
			CloudAccountID:       cloudAccount.ID,
			OrganizationalUnitID: ouID,
		})
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = tx.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, resp.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
