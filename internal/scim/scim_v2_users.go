package scim

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/scim/middleware"
	"github.com/suse-skyscraper/skyscraper/internal/scim/payloads"
	"github.com/suse-skyscraper/skyscraper/internal/scim/responses"
)

func V2ListUsers(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filterString := r.URL.Query().Get("filter")
		filters, err := ParseFilter(filterString)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadFilter(err))
			return
		}

		pagination := paginate(r)

		var totalCount int64
		var users []db.User

		if len(filters) == 0 {
			totalCount, err = app.DB.GetUserCount(r.Context())
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			users, err = app.DB.GetUsers(r.Context(), db.GetUsersParams{
				Offset: pagination.Offset,
				Limit:  pagination.Limit,
			})
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}
		} else {
			// we only support the userName filter for now
			// Okta uses this to see if a userName already exists
			filter := filters[0]
			if filter.FilterField == Username && filter.FilterOperator == eq {
				user, err := app.DB.FindByUsername(r.Context(), filter.FilterValue)
				switch err {
				case nil:
					users = []db.User{user}
					totalCount = int64(len(users))
				case pgx.ErrNoRows:
					users = []db.User{}
					totalCount = int64(len(users))
				default:
					_ = render.Render(w, r, responses.ErrInternalServerError)
					return
				}
			} else {
				_ = render.Render(w, r,
					responses.ErrBadFilter(
						errors.New("unsupported filter - only \"userName\" with the operator \"eq\" is supported"),
					),
				)
				return
			}
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserListResponse(
			app.Config,
			users,
			responses.ScimUserListResponseInput{
				StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
				TotalResults: int(totalCount),
				ItemsPerPage: int(pagination.Limit),
			}))
	}
}

func V2GetUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(app.Config, user))
	}
}

func V2CreateUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := payloads.UserPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		if payload.Username == "" {
			_ = render.Render(w, r, responses.ErrBadValue(errors.New("Attribute 'userName' is required")))
			return
		}

		user, err := app.DB.CreateUser(r.Context(), db.CreateUserParams{
			Username: payload.Username,
			Name:     payload.GetJSONName(),
			Active:   payload.Active,
			Emails:   payload.GetJSONEmails(),
			Locale: sql.NullString{
				String: payload.Locale,
				Valid:  payload.Locale != "",
			},
			ExternalID: sql.NullString{
				String: payload.ExternalID,
				Valid:  payload.ExternalID != "",
			},
			DisplayName: sql.NullString{
				String: payload.DisplayName,
				Valid:  payload.DisplayName != "",
			},
		})
		if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			w.WriteHeader(http.StatusConflict)
			return
		} else if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusCreated, responses.NewScimUserResponse(app.Config, user))
	}
}

func V2DeleteUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		err := app.DB.DeleteUser(r.Context(), user.ID)
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

func V2UpdateUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads.UserPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}
		err = app.DB.UpdateUser(r.Context(), db.UpdateUserParams{
			ID:       user.ID,
			Username: payload.Username,
			Name:     payload.GetJSONName(),
			Active:   payload.Active,
			Emails:   payload.GetJSONEmails(),
			Locale: sql.NullString{
				String: payload.Locale,
				Valid:  payload.Locale != "",
			},
			DisplayName: sql.NullString{
				String: payload.DisplayName,
				Valid:  payload.DisplayName != "",
			},
			ExternalID: sql.NullString{
				String: payload.ExternalID,
				Valid:  payload.ExternalID != "",
			},
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		user, err = app.DB.GetUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(app.Config, user))
	}
}

func V2PatchUser(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.User).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		payload, err := payloads.UserPatchPayloadFromJSON(r.Body)
		if err != nil {
			_ = render.Render(w, r, responses.ErrBadValue(err))
			return
		}

		for _, op := range payload.Operations {
			switch op.Op {
			case "replace":
				err = app.DB.PatchUser(r.Context(), db.PatchUserParams{
					ID:     user.ID,
					Active: op.Value.Active,
				})
				if err != nil {
					_ = render.Render(w, r, responses.ErrInternalServerError)
					return
				}
			default:
				_ = render.Render(w, r, responses.ErrBadValue(errors.New("Unsupported operation")))
			}
		}

		user, err = app.DB.GetUser(r.Context(), user.ID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserResponse(app.Config, user))
	}
}
