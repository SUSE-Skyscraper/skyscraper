package payloads

import (
	"net/http"

	"github.com/jackc/pgtype"
)

//----------------------------------------------------------------------------------------------------------------------
// Create or Update Resources
//----------------------------------------------------------------------------------------------------------------------

// CreateOrUpdateResourcePayloadData is the data for the CreateOrUpdateResourcePayload.
type CreateOrUpdateResourcePayloadData struct {
	AccountName string            `json:"account_name"`
	Active      bool              `json:"active"`
	TagsCurrent map[string]string `json:"tags_current"`
	TagsDesired map[string]string `json:"tags_desired"`
	tagsDesired pgtype.JSONB
	tagsCurrent pgtype.JSONB
}

// CreateOrUpdateResourcePayload is the payload for creating a cloud account.
type CreateOrUpdateResourcePayload struct {
	Data CreateOrUpdateResourcePayloadData `json:"data"`
}

// Bind binds extra data from the payload CreateOrUpdateResourcePayload.
func (u *CreateOrUpdateResourcePayload) Bind(_ *http.Request) error {
	// bind the tags current JSON to the payload
	tagsCurrent := pgtype.JSONB{}
	if u.Data.TagsCurrent != nil && len(u.Data.TagsCurrent) > 0 {
		err := tagsCurrent.Set(u.Data.TagsCurrent)
		if err != nil {
			return err
		}
	} else {
		err := tagsCurrent.Set("{}")
		if err != nil {
			return err
		}
	}

	u.Data.tagsCurrent = tagsCurrent

	// bind the tags desired JSON to the payload
	tagsDesired := pgtype.JSONB{}
	if u.Data.TagsDesired != nil && len(u.Data.TagsDesired) > 0 {
		err := tagsDesired.Set(u.Data.TagsDesired)
		if err != nil {
			return err
		}
	} else {
		err := tagsDesired.Set("{}")
		if err != nil {
			return err
		}
	}
	u.Data.tagsDesired = tagsDesired

	return nil
}

// GetTagsCurrent returns the parsed JSONB for the tags.
func (u *CreateOrUpdateResourcePayloadData) GetTagsCurrent() pgtype.JSONB {
	return u.tagsCurrent
}

// GetTagsDesired returns the parsed JSONB for the tags.
func (u *CreateOrUpdateResourcePayloadData) GetTagsDesired() pgtype.JSONB {
	return u.tagsDesired
}
