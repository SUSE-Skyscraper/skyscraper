package server

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/server/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/server/responses"
)

func V1Tags(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := app.Repository.GetTags(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewTagsResponse(tags))
	}
}

func V1UpdateTag(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tag, ok := r.Context().Value(middleware.ContextTag).(db.Tag)
		if !ok {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		var payload payloads.UpdateTagPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		tag, err = app.Repository.UpdateTag(r.Context(), db.UpdateTagParams{
			ID:          tag.ID,
			DisplayName: payload.Data.DisplayName,
			Description: payload.Data.Description,
			Key:         payload.Data.Key,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewTagResponse(tag))
	}
}

func V1CreateTag(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload payloads.CreateTagPayload
		err := render.Bind(r, &payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		tag, err := app.Repository.CreateTag(r.Context(), db.CreateTagParams{
			DisplayName: payload.Data.DisplayName,
			Description: payload.Data.Description,
			Key:         payload.Data.Key,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses.NewTagResponse(tag))
	}
}
