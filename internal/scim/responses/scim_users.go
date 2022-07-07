package responses

import (
	"fmt"
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type ScimUserResponse struct {
	Schemas  []string    `json:"schemas,omitempty"`
	UserName string      `json:"userName"`
	ID       string      `json:"id"`
	Name     interface{} `json:"name,omitempty"`
	Emails   interface{} `json:"emails,omitempty"`
	Active   bool        `json:"active"`
}

func (rd *ScimUserResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type ScimListUsersResponse struct {
	Schemas      []string            `json:"schemas"`
	ItemsPerPage int                 `json:"itemsPerPage"`
	StartIndex   int                 `json:"startIndex"`
	TotalResults int                 `json:"totalResults"`
	Resources    []*ScimUserResponse `json:"Resources"`
}

func (rd *ScimListUsersResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewScimUserResponse(user db.User) *ScimUserResponse {
	return newScimUserResponse(user, true)
}

func newScimUserResponse(user db.User, singleResponse bool) *ScimUserResponse {
	// the schemas should be added if the response is a single user, not a list
	var schemas []string
	if singleResponse {
		schemas = []string{"urn:ietf:params:scim:schemas:core:2.0:User"}
	}

	return &ScimUserResponse{
		Schemas:  schemas,
		ID:       fmt.Sprintf("%d", user.ID),
		UserName: user.Username,
		Name:     user.Name.Get(),
		Emails:   user.Emails.Get(),
		Active:   user.Active,
	}
}

type ScimUserListResponseInput struct {
	TotalResults int
	StartIndex   int
	ItemsPerPage int
}

func NewScimUserListResponse(users []db.User, input ScimUserListResponseInput) *ScimListUsersResponse {
	var list []*ScimUserResponse
	for _, user := range users {
		list = append(list, newScimUserResponse(user, false))
	}

	return &ScimListUsersResponse{
		Schemas:      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
		Resources:    list,
		TotalResults: input.TotalResults,
		StartIndex:   input.StartIndex,
		ItemsPerPage: input.ItemsPerPage,
	}
}
