package server

import (
	"context"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	db2 "github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	responses2 "github.com/suse-skyscraper/skyscraper/cli/internal/server/responses"

	"github.com/go-chi/render"
)

func V1StandardTags(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := app.Repository.GetTags(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses2.NewTagsResponse(tags))
	}
}

func V1UpdateStandardTag(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// bind the payload
		var payload payloads.UpdateTagPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		// Get the tag we want to change from the context
		tag, ok := r.Context().Value(middleware.ContextTag).(db2.StandardTag)
		if !ok {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		// Begin a database transaction
		repo, err := app.Repository.Begin(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// If any error occurs, rollback the transaction
		defer func(repo db2.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		// update the tag
		tag, err = repo.UpdateTag(r.Context(), db2.UpdateTagParams{
			ID:          tag.ID,
			DisplayName: payload.Data.DisplayName,
			Description: payload.Data.Description,
			Key:         payload.Data.Key,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// audit the change
		err = auditor.AuditChange(r.Context(), db2.AuditResourceTypeTag, tag.ID, payload)
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

		_ = render.Render(w, r, responses2.NewTagResponse(tag))
	}
}

func V1CreateStandardTag(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// bind the payload
		var payload payloads.CreateTagPayload
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
		defer func(repo db2.RepositoryQueries, ctx context.Context) {
			_ = repo.Rollback(ctx)
		}(repo, r.Context())

		// create an auditor within our transaction
		auditor := auditors.NewAuditor(repo)

		// create the tag
		tag, err := repo.CreateTag(r.Context(), db2.CreateTagParams{
			DisplayName: payload.Data.DisplayName,
			Description: payload.Data.Description,
			Key:         payload.Data.Key,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// audit the change
		err = auditor.AuditChange(r.Context(), db2.AuditResourceTypeTag, tag.ID, payload)
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

		_ = render.Render(w, r, responses2.NewTagResponse(tag))
	}
}
