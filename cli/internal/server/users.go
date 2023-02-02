package server

import (
	"net/http"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	db2 "github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/pagination"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"
	responses2 "github.com/suse-skyscraper/skyscraper/cli/internal/server/responses"

	"github.com/go-chi/render"
)

func V1Users(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		paginate := pagination.Paginate(r)
		users, err := app.Repository.GetUsers(r.Context(), db2.GetUsersParams{
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses2.NewUsersResponse(users))
	}
}

func V1User(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUser).(db2.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, responses2.NewUserResponse(user))
	}
}
