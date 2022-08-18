package responses

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

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

func NewTagResponse(tag db.StandardTag) *TagResponse {
	return &TagResponse{
		Data: newTagItem(tag),
	}
}

func NewTagsResponse(tags []db.StandardTag) *TagsResponse {
	list := make([]TagItem, len(tags))
	for i, tag := range tags {
		list[i] = newTagItem(tag)
	}

	return &TagsResponse{
		Data: list,
	}
}

func newTagItem(tag db.StandardTag) TagItem {
	return TagItem{
		ID:   tag.ID.String(),
		Type: ObjectResponseTypeTag,
		Attributes: TagItemAttributes{
			DisplayName: tag.DisplayName,
			Description: tag.Description,
			Key:         tag.Key,
			CreatedAt:   tag.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   tag.UpdatedAt.Format(time.RFC3339),
		},
	}
}
