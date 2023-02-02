package responses

import (
	"time"

	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

func NewTagResponse(tag db.StandardTag) *resp.TagResponse {
	return &resp.TagResponse{
		Data: newTagItem(tag),
	}
}

func NewTagsResponse(tags []db.StandardTag) *resp.TagsResponse {
	list := make([]resp.TagItem, len(tags))
	for i, tag := range tags {
		list[i] = newTagItem(tag)
	}

	return &resp.TagsResponse{
		Data: list,
	}
}

func newTagItem(tag db.StandardTag) resp.TagItem {
	return resp.TagItem{
		ID:   tag.ID.String(),
		Type: resp.ObjectResponseTypeTag,
		Attributes: resp.TagItemAttributes{
			DisplayName: tag.DisplayName,
			Description: tag.Description,
			Key:         tag.Key,
			CreatedAt:   tag.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   tag.UpdatedAt.Format(time.RFC3339),
		},
	}
}
