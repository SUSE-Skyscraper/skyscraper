package payloads

import (
	"net/http"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type UpdateCloudAccountPayload struct {
	TagsDesired map[string]string `json:"tags_desired"`
	json        pgtype.JSONB
}

func (u *UpdateCloudAccountPayload) Bind(_ *http.Request) error {
	if u.TagsDesired == nil {
		return errors.Errorf("tags_desired is required")
	}

	jsonTags := pgtype.JSONB{}
	err := jsonTags.Set(u.TagsDesired)
	if err != nil {
		return err
	}
	u.json = jsonTags

	return nil
}

func (u *UpdateCloudAccountPayload) GetJSON() pgtype.JSONB {
	return u.json
}
