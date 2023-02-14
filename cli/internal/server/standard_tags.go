package server

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/api/payloads"
	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/auditors"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
)

func V1StandardTags(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := app.Repo.GetTags(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewTagsResponse(tags))
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

		// Get the tag we want to change from the account
		tag, ok := r.Context().Value(middleware.ContextTag).(db.StandardTag)
		if !ok {
			_ = render.Render(w, r, responses.ErrNotFound)
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

		// update the tag
		tag, err = repo.UpdateTag(r.Context(), db.UpdateTagParams{
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
		err = auditor.AuditChange(r.Context(), db.AuditResourceTypeTag, tag.ID, payload)
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

		_ = render.Render(w, r, NewTagResponse(tag))
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

		// create the tag
		tag, err := repo.CreateTag(r.Context(), db.CreateTagParams{
			DisplayName: payload.Data.DisplayName,
			Description: payload.Data.Description,
			Key:         payload.Data.Key,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		// audit the change
		err = auditor.AuditChange(r.Context(), db.AuditResourceTypeTag, tag.ID, payload)
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

		_ = render.Render(w, r, NewTagResponse(tag))
	}
}

func NewTagResponse(tag db.StandardTag) *responses.TagResponse {
	return &responses.TagResponse{
		Data: newTagItem(tag),
	}
}

func NewTagsResponse(tags []db.StandardTag) *responses.TagsResponse {
	list := make([]responses.TagItem, len(tags))
	for i, tag := range tags {
		list[i] = newTagItem(tag)
	}

	return &responses.TagsResponse{
		Data: list,
	}
}

func newTagItem(tag db.StandardTag) responses.TagItem {
	return responses.TagItem{
		ID:   tag.ID.String(),
		Type: responses.ObjectResponseTypeTag,
		Attributes: responses.TagItemAttributes{
			DisplayName: tag.DisplayName,
			Description: tag.Description,
			Key:         tag.Key,
			CreatedAt:   tag.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   tag.UpdatedAt.Format(time.RFC3339),
		},
	}
}
