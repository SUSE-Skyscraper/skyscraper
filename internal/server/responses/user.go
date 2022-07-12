package responses

import (
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type UserAttributes struct {
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Active    bool   `json:"active"`
	Locale    string `json:"locale"`
}

type UserItem struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Attributes UserAttributes `json:"attributes"`
}

type UserResponse struct {
	Data UserItem `json:"data"`
}

func (rd *UserResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewUserResponse(user db.User) *UserResponse {
	return &UserResponse{
		Data: UserItem{
			ID:   user.ID.String(),
			Type: "user",
			Attributes: UserAttributes{
				Username:  user.Username,
				Active:    user.Active,
				Locale:    user.Locale.String,
				CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
			},
		},
	}
}
