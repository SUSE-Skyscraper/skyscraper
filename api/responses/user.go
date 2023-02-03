package responses

import "net/http"

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

type UsersResponse struct {
	Data []UserItem `json:"data"`
}

func (rd *UsersResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
