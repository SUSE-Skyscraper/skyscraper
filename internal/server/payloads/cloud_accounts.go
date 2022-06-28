package payloads

import (
	"net/http"

	"github.com/pkg/errors"
)

type UpdateCloudAccountPayload struct {
	TagsDesired map[string]string `json:"tags_desired"`
}

func (u *UpdateCloudAccountPayload) Bind(_ *http.Request) error {
	if u.TagsDesired == nil {
		return errors.Errorf("tags_desired is required")
	}

	return nil
}
