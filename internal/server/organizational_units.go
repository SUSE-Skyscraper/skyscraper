package server

import (
	"context"
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1ListOrganizationalUnits(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationalUnits, err := app.Repository.GetOrganizationalUnits(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewOrganizationalUnitsResponse(organizationalUnits))
	}
}

func V1GetOrganizationalUnit(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationalUnit, ok := r.Context().Value(middleware.ContextOrganizationalUnit).(db.OrganizationalUnit)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewOrganizationalUnitResponse(organizationalUnit))
	}
}

func V1DeleteOrganizationalUnit(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationalUnit, ok := r.Context().Value(middleware.ContextOrganizationalUnit).(db.OrganizationalUnit)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		children, err := repo.GetOrganizationalUnitChildren(r.Context(), organizationalUnit.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		} else if len(children) > 0 {
			_ = render.Render(w, r, responses.ErrOrganizationNotEmpty)
			return
		}

		accounts, err := repo.GetOrganizationalUnitCloudAccounts(r.Context(), organizationalUnit.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		} else if len(accounts) > 0 {
			_ = render.Render(w, r, responses.ErrOrganizationNotEmpty)
			return
		}

		err = repo.DeleteOrganizationalUnit(r.Context(), organizationalUnit.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// audit the change
		err = auditor.AuditDelete(r.Context(), db.AuditResourceTypeOrganizationalUnit, organizationalUnit.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = app.FGAClient.RemoveOrganizationalUnitRelationships(r.Context(), organizationalUnit.ID, organizationalUnit.ParentID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func V1CreateOrganizationalUnit(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// bind the payload
		var payload payloads.CreateOrganizationalUnitsPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		organizationalUnit, err := repo.CreateOrganizationalUnit(r.Context(), db.CreateOrganizationalUnitParams{
			DisplayName: payload.Data.DisplayName,
			ParentID:    payload.Data.GetParentID(),
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = app.FGAClient.AddOrganizationalUnit(r.Context(), organizationalUnit.ID, payload.Data.GetParentID())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// audit the change
		err = auditor.AuditCreate(r.Context(), db.AuditResourceTypeOrganizationalUnit, organizationalUnit.ID, payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// Commit the transaction
		err = repo.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		_ = render.Render(w, r, responses.NewOrganizationalUnitResponse(organizationalUnit))
	}
}
