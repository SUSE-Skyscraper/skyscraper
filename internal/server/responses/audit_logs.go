package responses

import (
	"net/http"
	"reflect"
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type AuditLogItemAttributes struct {
	Message      string               `json:"message"`
	CallerID     string               `json:"caller_id"`
	CallerType   string               `json:"caller_type"`
	ResourceType db.AuditResourceType `json:"resource_type"`
	ResourceID   string               `json:"resource_id"`
	CreatedAt    string               `json:"created_at"`
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

func NewAuditLogsListResponse(logs []db.AuditLog, callers []any) *AuditLogsResponse {
	logList := make([]AuditLogItem, len(logs))
	for i, log := range logs {
		logList[i] = newAuditLogItem(log)
	}

	includedList := make([]any, len(callers))
	for i, caller := range callers {
		switch reflect.TypeOf(caller).String() {
		case "db.User":
			includedList[i] = newUserItem(caller.(db.User))
		case "db.ApiKey":
			includedList[i] = newAPIKeyItem(caller.(db.ApiKey))
		}
	}

	return &AuditLogsResponse{
		Data:     logList,
		Included: includedList,
	}
}

func newAuditLogItem(log db.AuditLog) AuditLogItem {
	return AuditLogItem{
		ID:   log.ID.String(),
		Type: ObjectResponseTypeAuditLog,
		Attributes: AuditLogItemAttributes{
			Message:      log.Message,
			CallerID:     log.CallerID.String(),
			CallerType:   string(log.CallerType),
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID.String(),
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		},
		Relationships: map[ObjectResponseType]Relationship{
			ObjectResponseTypeUser: {
				RelationshipData: RelationshipData{
					ID:   log.CallerID.String(),
					Type: string(parseDBCallerType(log.CallerType)),
				},
			},
		},
	}
}
