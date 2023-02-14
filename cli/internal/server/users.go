package server

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/pagination"
	"github.com/suse-skyscraper/skyscraper/cli/internal/server/middleware"

	"github.com/go-chi/render"
)

func V1Users(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		paginate := pagination.Paginate(r)
		users, err := app.Repo.GetUsers(r.Context(), db.GetUsersParams{
			Limit:  paginate.Limit,
			Offset: paginate.Offset,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewUsersResponse(users))
	}
}

func V1User(_ *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(middleware.ContextUser).(db.User)
		if !ok {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewUserResponse(user))
	}
}

func NewUserResponse(user db.User) *responses.UserResponse {
	return &responses.UserResponse{
		Data: newUserItem(user),
	}
}

func NewUsersResponse(users []db.User) *responses.UsersResponse {
	list := make([]responses.UserItem, len(users))
	for i, user := range users {
		list[i] = newUserItem(user)
	}

	return &responses.UsersResponse{
		Data: list,
	}
}

func newUserItem(user db.User) responses.UserItem {
	return responses.UserItem{
		ID:   user.ID.String(),
		Type: "user",
		Attributes: responses.UserAttributes{
			Username:  user.Username,
			Active:    user.Active,
			Locale:    user.Locale.String,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}
}
