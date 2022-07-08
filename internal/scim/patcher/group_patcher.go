package patcher

import (
	"context"
	"errors"

	"github.com/google/uuid"
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

func (p *GroupPatcher) patchAdd(q *db.Queries, group db.Group, op *payloads.GroupPatchOperation) error {
	newMembers, err := op.GetAddMembersPatch()
	if err != nil {
		return errors.New("failed to get add members patch")
	}

	err = p.addMembers(q, group, newMembers)
	if err != nil {
		return errors.New("failed to add members")
	}

	return nil
}

func (p *GroupPatcher) patchRemove(q *db.Queries, group db.Group, op *payloads.GroupPatchOperation) error {
	id, err := op.ParseIDFromPath()
	if err != nil {
		return errors.New("failed to parse id from path")
	}

	return removeUserFromGroup(p.ctx, q, id, group.ID)
}

func (p *GroupPatcher) patchReplace(q *db.Queries, group db.Group, op *payloads.GroupPatchOperation) error {
	switch op.Path {
	case "members":
		newMembers, err := op.GetAddMembersPatch()
		if err != nil {
			return errors.New("failed to get add members patch")
		}

		err = p.replaceMembers(q, group, newMembers)
		if err != nil {
			return errors.New("failed to replace members")
		}
	default:
		patch, err := op.GetPatch()
		if err != nil {
			return errors.New("failed to get patch")
		}

		err = q.PatchGroupDisplayName(p.ctx, db.PatchGroupDisplayNameParams{
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
	tx, err := p.app.PostgresPool.Begin(p.ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, p.ctx)

	q := db.New(tx)

	for _, op := range payload.Operations {
		switch op.Op {
		case "add":
			err = p.patchAdd(q, group, op)
			if err != nil {
				return err
			}
		case "remove":
			err = p.patchRemove(q, group, op)
			if err != nil {
				return err
			}
		case "replace":
			err = p.patchReplace(q, group, op)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown operation")
		}
	}

	return tx.Commit(p.ctx)
}

func (p *GroupPatcher) replaceMembers(q *db.Queries, group db.Group, members []payloads.MemberPatch) error {
	err := q.DropMembershipForGroup(p.ctx, group.ID)
	if err != nil {
		return errors.New("failed to drop membership")
	}

	err = q.RemovePoliciesForGroup(p.ctx, group.ID.String())
	if err != nil {
		return errors.New("failed to drop membership permissions")
	}

	for _, member := range members {
		err = addUserToGroup(p.ctx, q, member.Value, group.ID)
		if err != nil {
			return errors.New("failed to add membership")
		}
	}

	return nil
}

func (p *GroupPatcher) addMembers(q *db.Queries, group db.Group, members []payloads.MemberPatch) error {
	for _, member := range members {
		err := addUserToGroup(p.ctx, q, member.Value, group.ID)
		if err != nil {
			return errors.New("failed to add membership")
		}
	}

	return nil
}

func addUserToGroup(ctx context.Context, q *db.Queries, userID, groupID uuid.UUID) error {
	err := q.CreateMembershipForUserAndGroup(ctx, db.CreateMembershipForUserAndGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	membership, err := q.GetGroupMembershipForUser(ctx, db.GetGroupMembershipForUserParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	} else if !membership.Username.Valid {
		return errors.New("internal error")
	}

	err = q.AddPolicy(ctx, db.AddPolicyParams{
		Ptype: "g",
		V0:    membership.Username.String,
		V1:    groupID.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func removeUserFromGroup(ctx context.Context, q *db.Queries, userID, groupID uuid.UUID) error {
	membership, err := q.GetGroupMembershipForUser(ctx, db.GetGroupMembershipForUserParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	} else if !membership.Username.Valid {
		return errors.New("internal error")
	}

	err = q.DropMembershipForUserAndGroup(ctx, db.DropMembershipForUserAndGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return errors.New("failed to drop membership")
	}

	err = q.RemovePolicy(ctx, db.RemovePolicyParams{
		Ptype: "g",
		V0:    membership.Username.String,
		V1:    groupID.String(),
	})
	if err != nil {
		return err
	}

	return nil
}
