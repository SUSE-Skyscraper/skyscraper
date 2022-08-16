package auditors

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/suse-skyscraper/skyscraper/internal/auth"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/server/middleware"
)

type Auditor struct {
	repo db.RepositoryQueries
}

func NewAuditor(repo db.RepositoryQueries) Auditor {
	return Auditor{
		repo: repo,
	}
}

func (a *Auditor) AuditDelete(ctx context.Context, resourceType db.AuditResourceType, resourceID uuid.UUID) error {
	message := fmt.Sprintf("deleted %s", resourceType)

	return a.audit(ctx, resourceType, resourceID, message)
}

func (a *Auditor) AuditCreate(ctx context.Context, resourceType db.AuditResourceType, resourceID uuid.UUID, payload any) error {
	jsonState, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("created with payload: %s", string(jsonState))
	return a.audit(ctx, resourceType, resourceID, message)
}

func (a *Auditor) AuditChange(ctx context.Context, resourceType db.AuditResourceType, resourceID uuid.UUID, payload any) error {
	jsonState, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("action with payload: %s", string(jsonState))

	return a.audit(ctx, resourceType, resourceID, message)
}

func (a *Auditor) audit(ctx context.Context, resourceType db.AuditResourceType, resourceID uuid.UUID, message string) error {
	caller, ok := ctx.Value(middleware.ContextCaller).(auth.Caller)
	if !ok {
		return errors.New("failed to get caller from context")
	}

	callerType, err := caller.GetDBType()
	if err != nil {
		return err
	}

	_, err = a.repo.CreateAuditLog(ctx, db.CreateAuditLogParams{
		CallerID:     caller.ID,
		CallerType:   callerType,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Message:      message,
	})
	if err != nil {
		return err
	}

	return nil
}