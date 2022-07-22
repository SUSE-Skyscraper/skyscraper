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

func (a *Auditor) Audit(ctx context.Context, resourceType db.AuditResourceType, resourceID uuid.UUID, state any) error {
	caller, ok := ctx.Value(middleware.ContextCaller).(auth.Caller)
	if !ok {
		return errors.New("failed to get caller from context")
	}

	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
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
		Message:      fmt.Sprintf("action with payload: %s", string(jsonState)),
	})
	if err != nil {
		return err
	}

	return nil
}
