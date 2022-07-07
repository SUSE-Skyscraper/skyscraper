package patcher

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/scim/payloads"
)

type GroupPatcher struct {
	ctx context.Context
	app *application.App
}

func NewGroupPatcher(ctx context.Context, app *application.App) *GroupPatcher {
	return &GroupPatcher{
		ctx: ctx,
		app: app,
	}
}

func (p *GroupPatcher) patchAdd(group db.Group, op *payloads.GroupPatchOperation) error {
	newMembers, err := op.GetAddMembersPatch()
	if err != nil {
		return errors.New("failed to get add members patch")
	}

	err = p.addMembers(group, newMembers)
	if err != nil {
		return errors.New("failed to add members")
	}

	return nil
}

func (p *GroupPatcher) patchRemove(group db.Group, op *payloads.GroupPatchOperation) error {
	id, err := op.ParseIDFromPath()
	if err != nil {
		return errors.New("failed to parse id from path")
	}

	err = p.app.DB.DropMembershipForUserAndGroup(p.ctx, db.DropMembershipForUserAndGroupParams{
		UserID:  id,
		GroupID: group.ID,
	})
	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return errors.New("failed to drop membership")
	}

	return nil
}

func (p *GroupPatcher) patchReplace(group db.Group, op *payloads.GroupPatchOperation) error {
	switch op.Path {
	case "members":
		newMembers, err := op.GetAddMembersPatch()
		if err != nil {
			return errors.New("failed to get add members patch")
		}

		err = p.replaceMembers(group, newMembers)
		if err != nil {
			return errors.New("failed to replace members")
		}
	default:
		patch, err := op.GetPatch()
		if err != nil {
			return errors.New("failed to get patch")
		}

		err = p.app.DB.PatchGroupDisplayName(p.ctx, db.PatchGroupDisplayNameParams{
			ID:          group.ID,
			DisplayName: patch.DisplayName,
		})
		if err != nil {
			return errors.New("failed to patch display name")
		}
	}

	return nil
}

func (p *GroupPatcher) Patch(group db.Group, payload *payloads.GroupPatchPayload) error {
	for _, op := range payload.Operations {
		switch op.Op {
		case "add":
			return p.patchAdd(group, op)
		case "remove":
			return p.patchRemove(group, op)
		case "replace":
			return p.patchReplace(group, op)
		default:
			return errors.New("unknown operation")
		}
	}

	return nil
}

func (p *GroupPatcher) replaceMembers(group db.Group, members []payloads.MemberPatch) error {
	tx, err := p.app.PostgresPool.Begin(p.ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, p.ctx)

	q := db.New(tx)

	err = q.DropMembershipForGroup(p.ctx, group.ID)
	if err != nil {
		return errors.New("failed to drop membership")
	}

	for _, member := range members {
		err = q.CreateMembershipForUserAndGroup(p.ctx, db.CreateMembershipForUserAndGroupParams{
			UserID:  member.Value,
			GroupID: group.ID,
		})
		if err != nil {
			return errors.New("failed to add membership")
		}
	}

	err = tx.Commit(p.ctx)
	if err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}

func (p *GroupPatcher) addMembers(group db.Group, members []payloads.MemberPatch) error {
	tx, err := p.app.PostgresPool.Begin(p.ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, p.ctx)

	q := db.New(tx)

	for _, member := range members {
		err = q.CreateMembershipForUserAndGroup(p.ctx, db.CreateMembershipForUserAndGroupParams{
			UserID:  member.Value,
			GroupID: group.ID,
		})
		if err != nil {
			return errors.New("failed to add membership")
		}
	}

	err = tx.Commit(p.ctx)
	if err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}
