package payloads

import (
	"net/http"

	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
)

type UpdateCloudAccountPayloadData struct {
	TagsDesired map[string]string `json:"tags_desired"`
	json        pgtype.JSONB
}
type UpdateCloudAccountPayload struct {
	Data UpdateCloudAccountPayloadData `json:"data"`
}

func (u *UpdateCloudAccountPayload) Bind(_ *http.Request) error {
	if u.Data.TagsDesired == nil {
		return errors.Errorf("tags_desired is required")
	}

	jsonTags := pgtype.JSONB{}
	err := jsonTags.Set(u.Data.TagsDesired)
	if err != nil {
		return err
	}
	u.Data.json = jsonTags

	return nil
}

func (u *UpdateCloudAccountPayloadData) GetJSON() pgtype.JSONB {
	return u.json
}
