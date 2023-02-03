package server

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/suse-skyscraper/skyscraper/api/responses"
	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

func V1ListAuditLogs(app *application.App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID := r.URL.Query().Get("resource_id")
		targetType := r.URL.Query().Get("resource_type")

		if targetID == "" || targetType == "" {
			logs, err := app.Repo.GetAuditLogs(r.Context())
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			callers, err := getCallersForLogs(r.Context(), app, logs)
			if err != nil {
				_ = render.Render(w, r, responses.ErrInternalServerError)
				return
			}

			_ = render.Render(w, r, NewAuditLogsListResponse(logs, callers))
			return
		}

		target := db.AuditResourceType(targetType)
		id, err := uuid.Parse(targetID)
		if err != nil {
			_ = render.Render(w, r, responses.ErrNotFound)
			return
		}

		logs, err := app.Repo.GetAuditLogsForTarget(r.Context(), db.GetAuditLogsForTargetParams{
			ResourceID:   id,
			ResourceType: target,
		})
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		callers, err := getCallersForLogs(r.Context(), app, logs)
		if err != nil {
			_ = render.Render(w, r, responses.ErrInternalServerError)
			return
		}

		_ = render.Render(w, r, NewAuditLogsListResponse(logs, callers))
	}
}

func getCallersForLogs(ctx context.Context, app *application.App, logs []db.AuditLog) ([]any, error) {
	checker := make(map[uuid.UUID]bool)
	var userIds []uuid.UUID
	var apiKeyIds []uuid.UUID

	for _, log := range logs {
		if checker[log.CallerID] {
			continue
		}
		checker[log.CallerID] = true
		switch log.CallerType {
		case db.CallerTypeUser:
			userIds = append(userIds, log.CallerID)
		case db.CallerTypeApiKey:
			apiKeyIds = append(apiKeyIds, log.CallerID)
		}
		userIds = append(userIds, log.CallerID)
	}

	users, err := app.Repo.GetUsersByID(ctx, userIds)
	if err != nil {
		return nil, err
	}

	apiKeys, err := app.Repo.FindAPIKeysByID(ctx, apiKeyIds)
	if err != nil {
		return nil, err
	}

	userCount := len(users)
	apiKeyCount := len(apiKeys)
	callers := make([]any, userCount+apiKeyCount)
	for i, user := range users {
		callers[i] = user
	}
	for i, apiKey := range apiKeys {
		callers[i+userCount] = apiKey
	}

	return callers, err
}

func NewAuditLogsListResponse(logs []db.AuditLog, callers []any) *responses.AuditLogsResponse {
	logList := make([]responses.AuditLogItem, len(logs))
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

	return &responses.AuditLogsResponse{
		Data:     logList,
		Included: includedList,
	}
}

func newAuditLogItem(log db.AuditLog) responses.AuditLogItem {
	return responses.AuditLogItem{
		ID:   log.ID.String(),
		Type: responses.ObjectResponseTypeAuditLog,
		Attributes: responses.AuditLogItemAttributes{
			Message:      log.Message,
			CallerID:     log.CallerID.String(),
			CallerType:   string(log.CallerType),
			ResourceType: string(log.ResourceType),
			ResourceID:   log.ResourceID.String(),
			CreatedAt:    log.CreatedAt.Format(time.RFC3339),
		},
		Relationships: map[responses.ObjectResponseType]responses.Relationship{
			responses.ObjectResponseTypeUser: {
				RelationshipData: responses.RelationshipData{
					ID:   log.CallerID.String(),
					Type: string(parseDBCallerType(log.CallerType)),
				},
			},
		},
	}
}

func parseDBCallerType(callerType db.CallerType) responses.ObjectResponseType {
	switch callerType {
	case db.CallerTypeUser:
		return responses.ObjectResponseTypeUser
	case db.CallerTypeApiKey:
		return responses.ObjectResponseTypeAPIKey
	default:
		return responses.ObjectResponseTypeUnknown
	}
}
