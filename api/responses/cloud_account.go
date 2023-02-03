package responses

import "net/http"

//----------------------------------------------------------------------------------------------------------------------
// Common
//----------------------------------------------------------------------------------------------------------------------

// CloudAccountItemAttributes represents the attributes of a cloud account.
type CloudAccountItemAttributes struct {
	CloudProvider     string            `json:"cloud_provider"`
	TenantID          string            `json:"tenant_id"`
	AccountID         string            `json:"account_id"`
	Name              string            `json:"name"`
	Active            bool              `json:"active"`
	TagsCurrent       map[string]string `json:"tags_current"`
	TagsDesired       map[string]string `json:"tags_desired"`
	TagsDriftDetected bool              `json:"tags_drift_detected"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
}

// CloudAccountItem represents a cloud account.
type CloudAccountItem struct {
	ID         string                     `json:"id"`
	Type       ObjectResponseType         `json:"type"`
	Attributes CloudAccountItemAttributes `json:"attributes"`
}

//----------------------------------------------------------------------------------------------------------------------
// Single
//----------------------------------------------------------------------------------------------------------------------

// CloudAccountResponse represents a single cloud account.
type CloudAccountResponse struct {
	Data CloudAccountItem `json:"data"`
}

// Render is a no-op.
func (rd *CloudAccountResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
// List
//----------------------------------------------------------------------------------------------------------------------

// CloudAccountListResponse represents a list of cloud accounts.
type CloudAccountListResponse struct {
	Data []CloudAccountItem `json:"data"`
}

// Render is a no-op.
func (rd *CloudAccountListResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
