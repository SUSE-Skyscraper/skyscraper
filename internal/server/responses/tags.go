package responses

import (
	"net/http"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type TagItemAttributes struct {
	DisplayName string                 `json:"display_name"`
	Required    bool                   `json:"required"`
	Description string                 `json:"description"`
	Key         string                 `json:"key"`
	Overrides   map[string]interface{} `json:"overrides"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

type TagItem struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes TagItemAttributes `json:"attributes"`
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

func NewTagResponse(tag db.Tag) *TagResponse {
	return &TagResponse{
		Data: newTagItem(tag),
	}
}

func NewTagsResponse(tags []db.Tag) *TagsResponse {
	list := make([]TagItem, len(tags))
	for i, tag := range tags {
		list[i] = newTagItem(tag)
	}

	return &TagsResponse{
		Data: list,
	}
}

func newTagItem(tag db.Tag) TagItem {
	return TagItem{
		ID:   tag.ID.String(),
		Type: "tag",
		Attributes: TagItemAttributes{
			DisplayName: tag.DisplayName,
			Required:    tag.Required,
			Description: tag.Description,
			Key:         tag.Key,
			CreatedAt:   tag.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   tag.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
}
