package db

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) getCallersForLogs(ctx context.Context, logs []AuditLog) ([]any, error) {
	checker := make(map[uuid.UUID]bool)
	var userIds []uuid.UUID
	var apiKeyIds []uuid.UUID

	for _, log := range logs {
		if checker[log.CallerID] {
			continue
		}
		checker[log.CallerID] = true
		switch log.CallerType {
		case CallerTypeUser:
			userIds = append(userIds, log.CallerID)
		case CallerTypeApiKey:
			apiKeyIds = append(apiKeyIds, log.CallerID)
		}
		userIds = append(userIds, log.CallerID)
	}

	users, err := r.db.GetUsersById(ctx, userIds)
	if err != nil {
		return nil, err
	}

	apiKeys, err := r.db.FindAPIKeysById(ctx, apiKeyIds)
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
