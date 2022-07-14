package db

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) getUsersForLogs(ctx context.Context, logs []AuditLog) ([]User, error) {
	checker := make(map[uuid.UUID]bool)
	var userIds []uuid.UUID

	for _, log := range logs {
		if checker[log.UserID] {
			continue
		}
		checker[log.UserID] = true
		userIds = append(userIds, log.UserID)
	}

	users, err := r.db.GetUsersById(ctx, userIds)
	if err != nil {
		return nil, err
	}

	return users, err
}
