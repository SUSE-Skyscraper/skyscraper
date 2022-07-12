package db

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/suse-skyscraper/skyscraper/internal/scim/filters"
	"github.com/suse-skyscraper/skyscraper/internal/scim/payloads"
)

var ErrConflict = errors.New("duplicate key value violates unique constraint")

var _ RepositoryQueries = (*Repository)(nil)

type Repository struct {
	search       Searcher
	postgresPool *pgxpool.Pool
	db           Querier
	tx           pgx.Tx
}

func (r *Repository) CreateCloudTenant(ctx context.Context, input CreateCloudTenantParams) error {
	return r.db.CreateCloudTenant(ctx, input)
}

func (r *Repository) UpdateCloudAccountTagsDriftDetected(
	ctx context.Context,
	input UpdateCloudAccountTagsDriftDetectedParams,
) error {
	return r.db.UpdateCloudAccountTagsDriftDetected(ctx, input)
}

func (r *Repository) CreateOrInsertCloudAccount(
	ctx context.Context,
	input CreateOrInsertCloudAccountParams,
) (CloudAccount, error) {
	return r.db.CreateOrInsertCloudAccount(ctx, input)
}

func (r *Repository) FindAPIKey(ctx context.Context, token string) (ScimApiKey, error) {
	return r.db.FindAPIKey(ctx, token)
}

func (r *Repository) InsertAPIKey(ctx context.Context, token string) (ScimApiKey, error) {
	return r.db.InsertAPIKey(ctx, token)
}

func (r *Repository) RemovePolicy(ctx context.Context, input RemovePolicyParams) error {
	return r.db.RemovePolicy(ctx, input)
}

func (r *Repository) CreatePolicy(ctx context.Context, input AddPolicyParams) error {
	return r.db.AddPolicy(ctx, input)
}

func (r *Repository) TruncatePolicies(ctx context.Context) error {
	return r.db.TruncatePolicies(ctx)
}

func (r *Repository) GetPolicies(ctx context.Context) ([]Policy, error) {
	return r.db.GetPolicies(ctx)
}

func (r *Repository) ScimPatchUser(ctx context.Context, input PatchUserParams) error {
	return r.db.PatchUser(ctx, input)
}

func (r *Repository) UpdateUser(ctx context.Context, id uuid.UUID, input UpdateUserParams) (User, error) {
	err := r.db.UpdateUser(ctx, input)
	if err != nil {
		return User{}, err
	}

	user, err := r.db.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteUser(ctx, id)
}

func (r *Repository) CreateUser(ctx context.Context, input CreateUserParams) (User, error) {
	user, err := r.db.CreateUser(ctx, input)
	if err != nil && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return User{}, ErrConflict
	} else if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetScimUsers(ctx context.Context, input GetScimUsersInput) (int64, []User, error) {
	if len(input.Filters) == 0 {
		totalCount, err := r.db.GetUserCount(ctx)
		if err != nil {
			return 0, nil, err
		}

		users, err := r.db.GetUsers(ctx, GetUsersParams{
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
		user, err := r.db.FindByUsername(ctx, filter.FilterValue)
		switch err {
		case nil:
			return 1, []User{user}, nil
		case pgx.ErrNoRows:
			return 0, []User{}, nil
		default:
			return 0, nil, err
		}
	} else {
		return 0, nil, errors.New("unsupported filter")
	}
}

func (r *Repository) GetCloudTenants(ctx context.Context) ([]CloudTenant, error) {
	return r.db.GetCloudTenants(ctx)
}

func (r *Repository) CreateGroup(ctx context.Context, displayName string) (Group, error) {
	group, err := r.db.CreateGroup(ctx, displayName)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (r *Repository) GetGroupMembership(ctx context.Context, idString string) ([]GetGroupMembershipRow, error) {
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	return r.db.GetGroupMembership(ctx, id)
}

func (r *Repository) GetGroups(ctx context.Context, params GetGroupsParams) (int64, []Group, error) {
	totalCount, err := r.db.GetGroupCount(ctx)
	if err != nil {
		return 0, nil, err
	}

	groups, err := r.db.GetGroups(ctx, params)
	if err != nil {
		return 0, nil, err
	}

	return totalCount, groups, nil
}

func (r *Repository) DeleteGroup(ctx context.Context, idString string) error {
	id, err := uuid.Parse(idString)
	if err != nil {
		return err
	}

	err = r.db.DeleteGroup(ctx, id)
	if err != nil {
		return err
	}

	err = r.db.RemovePoliciesForGroup(ctx, idString)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateGroup(ctx context.Context, input PatchGroupDisplayNameParams) (Group, error) {
	err := r.db.PatchGroupDisplayName(ctx, input)
	if err != nil {
		return Group{}, err
	}

	return r.FindGroup(ctx, input.ID.String())
}

func (r *Repository) Rollback(ctx context.Context) error {
	if r.tx == nil {
		return errors.New("no transaction in progress")
	}

	return r.tx.Rollback(ctx)
}

func (r *Repository) Begin(ctx context.Context) (RepositoryQueries, error) {
	tx, err := r.postgresPool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	db := &Queries{db: tx}

	return &Repository{
		search:       r.search,
		postgresPool: r.postgresPool,
		db:           db,
		tx:           tx,
	}, nil
}

func (r *Repository) Commit(ctx context.Context) error {
	if r.tx == nil {
		return errors.New("no transaction in progress")
	}

	return r.tx.Commit(ctx)
}

func (r *Repository) FindUser(ctx context.Context, id string) (User, error) {
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}

	return r.db.GetUser(ctx, idParsed)
}

func (r *Repository) FindGroup(ctx context.Context, id string) (Group, error) {
	idParsed, err := uuid.Parse(id)
	if err != nil {
		return Group{}, err
	}

	return r.db.GetGroup(ctx, idParsed)
}

func (r *Repository) SearchCloudAccounts(ctx context.Context, input SearchCloudAccountsInput) ([]CloudAccount, error) {
	return r.search.SearchCloudAccounts(ctx, input)
}

func (r *Repository) FindCloudAccount(ctx context.Context, input FindCloudAccountInput) (CloudAccount, error) {
	if input.ID != "" {
		id, err := uuid.Parse(input.ID)
		if err != nil {
			return CloudAccount{}, err
		}

		return r.db.FindCloudAccount(ctx, id)
	}

	return r.db.FindCloudAccountByCloudAndTenant(ctx, FindCloudAccountByCloudAndTenantParams{
		Cloud:     input.Cloud,
		TenantID:  input.TenantID,
		AccountID: input.AccountID,
	})
}

func (r *Repository) UpdateCloudAccount(ctx context.Context, input UpdateCloudAccountParams) (CloudAccount, error) {
	err := r.db.UpdateCloudAccount(ctx, input)
	if err != nil {
		return CloudAccount{}, err
	}

	return r.db.FindCloudAccount(ctx, input.ID)
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (User, error) {
	return r.db.FindByUsername(ctx, username)
}

func (r *Repository) AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []payloads.MemberPatch) error {
	for _, member := range members {
		err := r.AddUserToGroup(ctx, member.Value, groupID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ReplaceUsersInGroup(ctx context.Context, groupID uuid.UUID, members []payloads.MemberPatch) error {
	err := r.db.DropMembershipForGroup(ctx, groupID)
	if err != nil {
		return err
	}

	err = r.db.RemovePoliciesForGroup(ctx, groupID.String())
	if err != nil {
		return err
	}

	for _, member := range members {
		err = r.AddUserToGroup(ctx, member.Value, groupID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	err := r.db.CreateMembershipForUserAndGroup(ctx, CreateMembershipForUserAndGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	membership, err := r.db.GetGroupMembershipForUser(ctx, GetGroupMembershipForUserParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	err = r.db.AddPolicy(ctx, AddPolicyParams{
		Ptype: "g",
		V0:    membership.Username.String,
		V1:    groupID.String(),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	membership, err := r.db.GetGroupMembershipForUser(ctx, GetGroupMembershipForUserParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	err = r.db.DropMembershipForUserAndGroup(ctx, DropMembershipForUserAndGroupParams{
		UserID:  userID,
		GroupID: groupID,
	})
	if err != nil {
		return err
	}

	err = r.db.RemovePolicy(ctx, RemovePolicyParams{
		Ptype: "g",
		V0:    membership.Username.String,
		V1:    groupID.String(),
	})
	if err != nil {
		return err
	}

	return nil
}
