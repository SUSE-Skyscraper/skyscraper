package scim

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/scim/responses"
)

func V2ListGroups(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination, err := paginate(r)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInvalidRequest(err))
			return
		}

		totalCount, err := app.DB.GetUserCount(r.Context())
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		users, err := app.DB.GetUsers(r.Context(), db.GetUsersParams{
			Offset: pagination.Offset,
			Limit:  pagination.Limit,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		RenderScimJSON(w, r, http.StatusOK, responses.NewScimUserListResponse(users, responses.ScimUserListResponseInput{
			StartIndex:   int(pagination.Offset)/int(pagination.Limit) + 1,
			TotalResults: int(totalCount),
			ItemsPerPage: int(pagination.Limit),
		}))
	}
}
