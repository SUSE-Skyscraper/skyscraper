package payloads

import (
	"net/http"
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
}

// CreateOrUpdateResourcePayload is the payload for creating a cloud account.
type CreateOrUpdateResourcePayload struct {
	Data CreateOrUpdateResourcePayloadData `json:"data"`
}

// Bind binds extra data from the payload CreateOrUpdateResourcePayload.
func (u *CreateOrUpdateResourcePayload) Bind(_ *http.Request) error {
	// if nil, initialize the map
	if u.Data.TagsCurrent == nil {
		u.Data.TagsCurrent = make(map[string]string)
	}

	// if nil, initialize the map
	if u.Data.TagsDesired == nil {
		u.Data.TagsDesired = make(map[string]string)
	}

	return nil
}
