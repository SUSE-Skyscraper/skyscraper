package scim

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/scim/middleware"
	patcher "github.com/suse-skyscraper/skyscraper/internal/scim/patcher"
	"github.com/suse-skyscraper/skyscraper/internal/scim/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/scim/responses"
)

func V2ListGroups(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := paginate(r)

		totalCount, err := app.DB.GetGroupCount(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		groups, err := app.DB.GetGroups(r.Context(), db.GetGroupsParams{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupListResponse(
			app.Config,
			groups,
			responses.ScimGroupListResponseInput{
				StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(pagination.Limit),
			}))
	}
}

func V2GetGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		group := r.Context().Value(middleware.Group).(db.Group)

		members, err := app.DB.GetGroupMembership(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupResponse(app.Config, group, members))
	}
}

func V2CreateGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.GroupPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		group, err := app.DB.CreateGroup(r.Context(), payload.DisplayName)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}
		var members []db.GetGroupMembershipRow

		RenderScimJSON(w, r, http.StatusCreated, responses.NewScimGroupResponse(app.Config, group, members))
	}
}

func V2PatchGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.GroupPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		groupPatcher := patcher.NewGroupPatcher(r.Context(), app)
		group := r.Context().Value(middleware.Group).(db.Group)
		err = groupPatcher.Patch(group, payload)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		group, err = app.DB.GetGroup(r.Context(), group.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		var members []db.GetGroupMembershipRow
		RenderScimJSON(w, r, http.StatusOK, responses.NewScimGroupResponse(app.Config, group, members))
	}
}

func V2DeleteGroup(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		group := r.Context().Value(middleware.Group).(db.Group)

		err := app.DB.DeleteGroup(r.Context(), group.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			_ = render.Render(w, r, responses.ErrNotFound(id))
			return
		} else if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
