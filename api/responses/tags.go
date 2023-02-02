package responses

import "net/http"

type TagItemAttributes struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Key         string `json:"key"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TagItem struct {
	ID         string             `json:"id"`
	Type       ObjectResponseType `json:"type"`
	Attributes TagItemAttributes  `json:"attributes"`
}

type TagResponse struct {
	Data TagItem `json:"data"`
}

func (rd *TagResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type TagsResponse struct {
	Data []TagItem `json:"data"`
}

func (rd *TagsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
