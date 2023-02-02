package responses

import (
	"time"

	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

func NewUserResponse(user db.User) *resp.UserResponse {
	return &resp.UserResponse{
		Data: newUserItem(user),
	}
}

func NewUsersResponse(users []db.User) *resp.UsersResponse {
	list := make([]resp.UserItem, len(users))
	for i, user := range users {
		list[i] = newUserItem(user)
	}

	return &resp.UsersResponse{
		Data: list,
	}
}

func newUserItem(user db.User) resp.UserItem {
	return resp.UserItem{
		ID:   user.ID.String(),
		Type: "user",
		Attributes: resp.UserAttributes{
			Username:  user.Username,
			Active:    user.Active,
			Locale:    user.Locale.String,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}
}
