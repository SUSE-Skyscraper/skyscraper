package responses

import (
	"reflect"
	"time"

	resp "github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
)

func NewAuditLogsListResponse(logs []db.AuditLog, callers []any) *resp.AuditLogsResponse {
	logList := make([]resp.AuditLogItem, len(logs))
	for i, log := range logs {
		logList[i] = newAuditLogItem(log)
	}

	includedList := make([]any, len(callers))
	for i, caller := range callers {
		switch reflect.TypeOf(caller).String() {
		case "db.User":
			includedList[i] = newUserItem(caller.(db.User))
		case "db.ApiKey":
			includedList[i] = newAPIKeyItem(caller.(db.ApiKey), "")
		}
	}

	return &resp.AuditLogsResponse{
		Data:     logList,
		Included: includedList,
	}
}

func newAuditLogItem(log db.AuditLog) resp.AuditLogItem {
	return resp.AuditLogItem{
		ID:   log.ID.String(),
		Type: resp.ObjectResponseTypeAuditLog,
		Attributes: resp.AuditLogItemAttributes{
			Message:      log.Message,
			CallerID:     log.CallerID.String(),
			CallerType:   string(log.CallerType),
			ResourceType: string(log.ResourceType),
			ResourceID:   log.ResourceID.String(),
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		},
		Relationships: map[resp.ObjectResponseType]resp.Relationship{
			resp.ObjectResponseTypeUser: {
				RelationshipData: resp.RelationshipData{
					ID:   log.CallerID.String(),
					Type: string(parseDBCallerType(log.CallerType)),
				},
			},
		},
	}
}

func parseDBCallerType(callerType db.CallerType) resp.ObjectResponseType {
	switch callerType {
	case db.CallerTypeUser:
		return resp.ObjectResponseTypeUser
	case db.CallerTypeApiKey:
		return resp.ObjectResponseTypeAPIKey
	default:
		return resp.ObjectResponseTypeUnknown
	}
}
