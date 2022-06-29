package workers

type ChangeTagsPayload struct {
	Cloud     string `json:"cloud"`
	TenantID  string `json:"tenant_id"`
	AccountID string `json:"account_id"`
}
