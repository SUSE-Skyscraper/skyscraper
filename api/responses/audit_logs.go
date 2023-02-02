package responses

import "net/http"

type AuditLogItemAttributes struct {
	Message      string `json:"message"`
	CallerID     string `json:"caller_id"`
	CallerType   string `json:"caller_type"`
	ResourceType string `json:"resource_type"`
	ResourceID   string `json:"resource_id"`
	CreatedAt    string `json:"created_at"`
}

type AuditLogItem struct {
	ID            string                              `json:"id"`
	Type          ObjectResponseType                  `json:"type"`
	Attributes    AuditLogItemAttributes              `json:"attributes"`
	Relationships map[ObjectResponseType]Relationship `json:"relationships"`
}

type AuditLogsResponse struct {
	Data     []AuditLogItem `json:"data"`
	Included []any          `json:"included"`
}

func (rd *AuditLogsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
