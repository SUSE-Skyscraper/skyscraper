package scimbridgedb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/suse-skyscraper/openfga-scim-bridge/v2/filters"

	"github.com/jackc/pgx/v4"

	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"

	"github.com/google/uuid"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/database"
	"github.com/suse-skyscraper/openfga-scim-bridge/v2/payloads"
)

var _ database.Bridge = (*DB)(nil)

type DB struct {
	app *application.App
}

func New(app *application.App) DB {
	return DB{
		app: app,
	}
}

func (d *DB) PatchGroup(ctx context.Context, groupID uuid.UUID, operations []payloads.GroupPatchOperation) error {
	tx, err := d.app.PostgresPool.Begin(ctx)
	if err != nil {
		return errors.New("failed to begin transaction")
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	repo := d.app.Repo.WithTx(tx)

	for _, op := range operations {
		switch op.Op {
		case "add":
			err = d.patchAdd(ctx, repo, groupID, op)
			if err != nil {
				return err
			}
		case "remove":
			err = d.patchRemove(ctx, repo, groupID, op)
			if err != nil {
				return err
			}
		case "replace":
			err = d.patchReplace(ctx, repo, groupID, op)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown operation")
		}
	}

	return tx.Commit(ctx)
}

func (d *DB) patchAdd(ctx context.Context, repo db.Repository, groupID uuid.UUID, op payloads.GroupPatchOperation) error {
	newMembers, err := op.GetAddMembersPatch()
	if err != nil {
		return errors.New("failed to get add members patch")
	}

	err = d.app.FGAClient.AddUsersToGroup(ctx, newMembers, groupID)
	if err != nil {
		return errors.New("failed to add members to FGA group")
	}

	for _, member := range newMembers {
		err := repo.CreateMembershipForUserAndGroup(ctx, db.CreateMembershipForUserAndGroupParams{
			UserID:  member,
			GroupID: groupID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) patchRemove(ctx context.Context, repo db.Repository, groupID uuid.UUID, op payloads.GroupPatchOperation) error {
	id, err := op.ParseIDFromPath()
	if err != nil {
		return errors.New("failed to parse id from path")
	}

	err = d.app.FGAClient.RemoveUserFromGroup(ctx, id, groupID)
	if err != nil {
		return err
	}
	err = repo.DropMembershipForUserAndGroup(ctx, db.DropMembershipForUserAndGroupParams{
		UserID:  id,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) patchReplace(ctx context.Context, repo db.Repository, groupID uuid.UUID, op payloads.GroupPatchOperation) error {
	switch op.Path {
	case "members":
		newMembers, err := op.GetAddMembersPatch()
		if err != nil {
			return errors.New("failed to get add members patch")
		}

		err = d.app.FGAClient.ReplaceUsersInGroup(ctx, newMembers, groupID)
		if err != nil {
			return errors.New("failed to replace members in FGA")
		}

		err = repo.DropMembershipForGroup(ctx, groupID)
		if err != nil {
			return err
		}

		for _, member := range newMembers {
			err = repo.CreateMembershipForUserAndGroup(ctx, db.CreateMembershipForUserAndGroupParams{
				UserID:  member,
				GroupID: groupID,
			})
			if err != nil {
				return err
			}
		}
		if err != nil {
			return errors.New("failed to replace members")
		}
	default:
		patch, err := op.GetPatch()
		if err != nil {
			return errors.New("failed to get patch")
		}

		err = repo.PatchGroupDisplayName(ctx, db.PatchGroupDisplayNameParams{
			ID:          groupID,
			DisplayName: patch.DisplayName,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	err := d.app.FGAClient.RemoveUsersInGroup(ctx, groupID)
	if err != nil {
		return err
	}

	err = d.app.Repo.DeleteGroup(ctx, groupID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateGroup(ctx context.Context, displayName string) (database.Group, error) {
	group, err := d.app.Repo.CreateGroup(ctx, displayName)
	if err != nil {
		return database.Group{}, err
	}

	scimGroup := toScimGroup(group)
	return scimGroup, nil
}

func (d *DB) GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]database.GroupMembership, error) {
	members, err := d.app.Repo.GetGroupMembership(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var groupMembers []database.GroupMembership
	for _, member := range members {
		groupMembers = append(groupMembers, database.GroupMembership{
			GroupID:  member.GroupID,
			Username: member.Username,
			UserID:   member.UserID,
		})
	}

	return groupMembers, nil
}

func (d *DB) FindGroup(ctx context.Context, groupID uuid.UUID) (database.Group, error) {
	group, err := d.app.Repo.GetGroup(ctx, groupID)
	if err != nil {
		return database.Group{}, err
	}

	scimGroup := toScimGroup(group)
	return scimGroup, nil
}

func (d *DB) GetGroups(ctx context.Context, limit int32, offset int32) (int64, []database.Group, error) {
	totalCount, err := d.app.Repo.GetGroupCount(ctx)
	if err != nil {
		return 0, nil, err
	}

	groups, err := d.app.Repo.GetGroups(ctx, db.GetGroupsParams{
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		return 0, nil, err
	}

	var scimGroups []database.Group
	for _, group := range groups {
		scimGroups = append(scimGroups, toScimGroup(group))
	}

	return totalCount, scimGroups, nil
}

func (d *DB) FindUser(ctx context.Context, userID uuid.UUID) (database.User, error) {
	user, err := d.app.Repo.GetUser(ctx, userID)
	if err != nil {
		return database.User{}, err
	}

	scimUser, err := toScimUser(user)
	if err != nil {
		return database.User{}, err
	}
	return scimUser, nil
}

func (d *DB) SetUserActive(ctx context.Context, userID uuid.UUID, active bool) error {
	err := d.app.Repo.PatchUser(ctx, db.PatchUserParams{
		ID:     userID,
		Active: active,
	})
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) UpdateUser(ctx context.Context, userID uuid.UUID, arg database.UserParams) (database.User, error) {
	name, err := parseJSONB(arg.Name)
	if err != nil {
		return database.User{}, err
	}

	emails, err := parseJSONB(arg.Emails)
	if err != nil {
		return database.User{}, err
	}

	user, err := d.app.Repo.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userID,
		Username: arg.Username,
		Name:     name,
		Active:   arg.Active,
		Emails:   emails,
		Locale: sql.NullString{
			String: arg.Locale,
			Valid:  arg.Locale != "",
		},
		DisplayName: sql.NullString{
			String: arg.DisplayName,
			Valid:  arg.DisplayName != "",
		},
		ExternalID: sql.NullString{
			String: arg.ExternalID,
			Valid:  arg.ExternalID != "",
		},
	})
	if err != nil {
		return database.User{}, err
	}

	scimUser, err := toScimUser(user)
	if err != nil {
		return database.User{}, err
	}

	return scimUser, nil
}

func (d *DB) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	tx, err := d.app.PostgresPool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	repo := d.app.Repo.WithTx(tx)

	err = repo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	err = d.app.FGAClient.RemoveUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateUser(ctx context.Context, arg database.UserParams) (database.User, error) {
	tx, err := d.app.PostgresPool.Begin(ctx)
	if err != nil {
		return database.User{}, err
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	repo := d.app.Repo.WithTx(tx)

	name, err := parseJSONB(arg.Name)
	if err != nil {
		return database.User{}, err
	}

	emails, err := parseJSONB(arg.Emails)
	if err != nil {
		return database.User{}, err
	}

	user, err := repo.CreateUser(ctx, db.CreateUserParams{
		Username: arg.Username,
		Name:     name,
		Active:   arg.Active,
		Emails:   emails,
		Locale: sql.NullString{
			String: arg.Locale,
			Valid:  arg.Locale != "",
		},
		ExternalID: sql.NullString{
			String: arg.ExternalID,
			Valid:  arg.ExternalID != "",
		},
		DisplayName: sql.NullString{
			String: arg.DisplayName,
			Valid:  arg.DisplayName != "",
		},
	})

	if err != nil {
		return database.User{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return database.User{}, err
	}

	scimUser, err := toScimUser(user)
	if err != nil {
		return database.User{}, err
	}

	return scimUser, nil
}

func (d *DB) GetUsers(ctx context.Context, input database.GetUsersParams) (int64, []database.User, error) {
	count, users, err := d.getUsers(ctx, input)
	if err != nil {
		return 0, nil, err
	}

	var scimUsers []database.User
	for _, user := range users {
		scimUser, err := toScimUser(user)
		if err != nil {
			return 0, nil, err
		}

		scimUsers = append(scimUsers, scimUser)
	}

	return count, scimUsers, nil
}

func (d *DB) getUsers(ctx context.Context, input database.GetUsersParams) (int64, []db.User, error) {
	if len(input.Filters) == 0 {
		totalCount, err := d.app.Repo.GetUserCount(ctx)
		if err != nil {
			return 0, nil, err
		}

		users, err := d.app.Repo.GetUsers(ctx, db.GetUsersParams{
			Offset: input.Offset,
			Limit:  input.Limit,
		})
		if err != nil {
			return 0, nil, err
		}

		return totalCount, users, nil
	}

	// we only support the userName filter for now
	// Okta uses this to see if a userName already exists
	filter := input.Filters[0]
	if filter.FilterField == filters.Username && filter.FilterOperator == filters.Eq {
		user, err := d.app.Repo.FindUserByUsername(ctx, filter.FilterValue)
		switch err {
		case nil:
			return 1, []db.User{user}, nil
		case pgx.ErrNoRows:
			return 0, []db.User{}, nil
		default:
			return 0, nil, err
		}
	} else {
		return 0, nil, errors.New("unsupported filter")
	}
}
