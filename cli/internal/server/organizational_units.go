package server

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4"

	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/render"
)

func V1ListOrganizationalUnits(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationalUnits, err := app.Repo.GetOrganizationalUnits(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewOrganizationalUnitsResponse(organizationalUnits))
	}
}

func V1GetOrganizationalUnit(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		organizationalUnit, ok := r.Context().Value(middleware.ContextOrganizationalUnit).(db.OrganizationalUnit)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewOrganizationalUnitResponse(organizationalUnit))
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
		tx, err := app.PostgresPool.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}
		repo := app.Repo.WithTx(tx)

		// If any error occurs, rollback the transaction
		defer func(repo pgx.Tx, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(tx, r.Context())

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
		err = tx.Commit(r.Context())
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
		tx, err := app.PostgresPool.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}
		repo := app.Repo.WithTx(tx)

		// If any error occurs, rollback the transaction
		defer func(tx pgx.Tx, ctx context.Context) {
			_ = tx.Rollback(ctx)
		}(tx, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)
		var parentID uuid.NullUUID
		if payload.Data.ParentID != "" {
			id, err := uuid.Parse(payload.Data.ParentID)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
			parentID = uuid.NullUUID{UUID: id, Valid: true}
		} else {
			parentID = uuid.NullUUID{UUID: uuid.Nil, Valid: false}
		}

		organizationalUnit, err := repo.CreateOrganizationalUnit(r.Context(), db.CreateOrganizationalUnitParams{
			DisplayName: payload.Data.DisplayName,
			ParentID:    parentID,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err = app.FGAClient.AddOrganizationalUnit(r.Context(), organizationalUnit.ID, parentID)
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
		err = tx.Commit(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		render.Status(r, http.StatusCreated)
		_ = render.Render(w, r, NewOrganizationalUnitResponse(organizationalUnit))
	}
}

func newOrganizationalUnitItem(organizationalUnit db.OrganizationalUnit) responses.OrganizationalUnitItem {
	parentID := ""
	if organizationalUnit.ParentID.Valid {
		parentID = organizationalUnit.ParentID.UUID.String()
	}

	return responses.OrganizationalUnitItem{
		ID:   organizationalUnit.ID.String(),
		Type: responses.ObjectResponseTypeOrganizationalUnit,
		Attributes: responses.OrganizationalUnitAttributes{
			ParentID:    parentID,
			DisplayName: organizationalUnit.DisplayName,
		},
	}
}

func NewOrganizationalUnitResponse(organizationalUnit db.OrganizationalUnit) *responses.OrganizationalUnitResponse {
	return &responses.OrganizationalUnitResponse{
		Data: newOrganizationalUnitItem(organizationalUnit),
	}
}

func NewOrganizationalUnitsResponse(organizationalUnits []db.OrganizationalUnit) *responses.OrganizationalUnitsResponse {
	list := make([]responses.OrganizationalUnitItem, len(organizationalUnits))
	for i, ou := range organizationalUnits {
		list[i] = newOrganizationalUnitItem(ou)
	}

	return &responses.OrganizationalUnitsResponse{
		Data: list,
	}
}
