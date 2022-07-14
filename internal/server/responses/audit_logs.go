package responses

import (
	"net/http"
	"time"

	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type AuditLogItemAttributes struct {
	Message      string               `json:"message"`
	UserID       string               `json:"user_id"`
	Username     string               `json:"username"`
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
	Included []UserItem     `json:"included"`
}

func (rd *AuditLogsResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewAuditLogsListResponse(logs []db.AuditLog, users []db.User) *AuditLogsResponse {
	logList := make([]AuditLogItem, len(logs))
	for i, log := range logs {
		logList[i] = newAuditLogItem(log)
	}

	userList := make([]UserItem, len(users))
	for i, user := range users {
		userList[i] = newUserItem(user)
	}

	return &AuditLogsResponse{
		Data:     logList,
		Included: userList,
	}
}

func newAuditLogItem(log db.AuditLog) AuditLogItem {
	return AuditLogItem{
		ID:   log.ID.String(),
		Type: ObjectResponseTypeAuditLog,
		Attributes: AuditLogItemAttributes{
			Message:      log.Message,
			UserID:       log.UserID.String(),
			ResourceType: log.ResourceType,
			ResourceID:   log.ResourceID.String(),
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		},
		Relationships: map[ObjectResponseType]Relationship{
			ObjectResponseTypeUser: {
				RelationshipData: RelationshipData{
					ID:   log.UserID.String(),
					Type: "user",
				},
			},
		},
	}
}
