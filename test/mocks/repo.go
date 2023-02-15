package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/cli/db"
)

type TestRepo struct {
	mock.Mock
}

func (t *TestRepo) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (db.User, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepo) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	args := t.Called(ctx, username)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepo) UpdateTag(ctx context.Context, arg db.UpdateTagParams) (db.StandardTag, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.StandardTag), args.Error(1)
}

func (t *TestRepo) GetOrganizationalUnitChildren(ctx context.Context, parentID uuid.UUID) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx, parentID)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepo) AssignAccountToOU(ctx context.Context, arg db.AssignAccountToOUParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) CreateAuditLog(ctx context.Context, arg db.CreateAuditLogParams) (db.AuditLog, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.AuditLog), args.Error(1)
}

func (t *TestRepo) CreateGroup(ctx context.Context, displayName string) (db.Group, error) {
	args := t.Called(ctx, displayName)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestRepo) CreateMembershipForUserAndGroup(ctx context.Context, arg db.CreateMembershipForUserAndGroupParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) CreateOrUpdateCloudAccount(ctx context.Context, arg db.CreateOrUpdateCloudAccountParams) (db.CloudAccount, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestRepo) CreateOrUpdateCloudTenant(ctx context.Context, arg db.CreateOrUpdateCloudTenantParams) (db.CloudTenant, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudTenant), args.Error(1)
}

func (t *TestRepo) CreateOrganizationalUnit(ctx context.Context, arg db.CreateOrganizationalUnitParams) (db.OrganizationalUnit, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepo) CreateTag(ctx context.Context, arg db.CreateTagParams) (db.StandardTag, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.StandardTag), args.Error(1)
}

func (t *TestRepo) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepo) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepo) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepo) DeleteOrganizationalUnit(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepo) DeleteScimAPIKey(ctx context.Context) error {
	args := t.Called(ctx)

	return args.Error(0)
}

func (t *TestRepo) DeleteTag(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepo) DropMembershipForGroup(ctx context.Context, groupID uuid.UUID) error {
	args := t.Called(ctx, groupID)

	return args.Error(0)
}

func (t *TestRepo) DropMembershipForUserAndGroup(ctx context.Context, arg db.DropMembershipForUserAndGroupParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) FindAPIKey(ctx context.Context, id uuid.UUID) (db.ApiKey, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepo) FindAPIKeysByID(ctx context.Context, id []uuid.UUID) ([]db.ApiKey, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.ApiKey), args.Error(1)
}

func (t *TestRepo) FindCloudAccount(ctx context.Context, id uuid.UUID) (db.CloudAccount, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestRepo) FindCloudAccountByCloudAndTenant(ctx context.Context, arg db.FindCloudAccountByCloudAndTenantParams) (db.CloudAccount, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestRepo) FindOrganizationalUnit(ctx context.Context, id uuid.UUID) (db.OrganizationalUnit, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepo) FindScimAPIKey(ctx context.Context) (db.ApiKey, error) {
	args := t.Called(ctx)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepo) FindTag(ctx context.Context, id uuid.UUID) (db.StandardTag, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.StandardTag), args.Error(1)
}

func (t *TestRepo) GetAPIKeys(ctx context.Context) ([]db.ApiKey, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.ApiKey), args.Error(1)
}

func (t *TestRepo) GetAPIKeysOrganizationalUnits(ctx context.Context, apiKeyID uuid.UUID) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx, apiKeyID)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepo) GetAuditLogs(ctx context.Context) ([]db.AuditLog, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.AuditLog), args.Error(1)
}

func (t *TestRepo) GetAuditLogsForTarget(ctx context.Context, arg db.GetAuditLogsForTargetParams) ([]db.AuditLog, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.AuditLog), args.Error(1)
}

func (t *TestRepo) GetCloudTenant(ctx context.Context, arg db.GetCloudTenantParams) (db.CloudTenant, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudTenant), args.Error(1)
}

func (t *TestRepo) GetCloudTenants(ctx context.Context) ([]db.CloudTenant, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.CloudTenant), args.Error(1)
}

func (t *TestRepo) GetGroup(ctx context.Context, id uuid.UUID) (db.Group, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestRepo) GetGroupCount(ctx context.Context) (int64, error) {
	args := t.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

func (t *TestRepo) GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]db.GetGroupMembershipRow, error) {
	args := t.Called(ctx, groupID)

	return args.Get(0).([]db.GetGroupMembershipRow), args.Error(1)
}

func (t *TestRepo) GetGroupMembershipForUser(ctx context.Context, arg db.GetGroupMembershipForUserParams) (db.GetGroupMembershipForUserRow, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.GetGroupMembershipForUserRow), args.Error(1)
}

func (t *TestRepo) GetGroups(ctx context.Context, arg db.GetGroupsParams) ([]db.Group, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.Group), args.Error(1)
}

func (t *TestRepo) GetOrganizationalUnitCloudAccounts(ctx context.Context, organizationalUnitID uuid.UUID) ([]db.CloudAccount, error) {
	args := t.Called(ctx, organizationalUnitID)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestRepo) GetOrganizationalUnits(ctx context.Context) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepo) GetTags(ctx context.Context) ([]db.StandardTag, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.StandardTag), args.Error(1)
}

func (t *TestRepo) GetUser(ctx context.Context, id uuid.UUID) (db.User, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepo) GetUserCount(ctx context.Context) (int64, error) {
	args := t.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

func (t *TestRepo) GetUserOrganizationalUnits(ctx context.Context, userID uuid.UUID) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx, userID)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepo) GetUsers(ctx context.Context, arg db.GetUsersParams) ([]db.User, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.User), args.Error(1)
}

func (t *TestRepo) GetUsersByID(ctx context.Context, userIDs []uuid.UUID) ([]db.User, error) {
	args := t.Called(ctx, userIDs)

	return args.Get(0).([]db.User), args.Error(1)
}

func (t *TestRepo) InsertAPIKey(ctx context.Context, arg db.InsertAPIKeyParams) (db.ApiKey, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepo) InsertScimAPIKey(ctx context.Context, apiKeyID uuid.UUID) (db.ScimApiKey, error) {
	args := t.Called(ctx, apiKeyID)

	return args.Get(0).(db.ScimApiKey), args.Error(1)
}

func (t *TestRepo) OrganizationalUnitsCloudAccounts(ctx context.Context, id []uuid.UUID) ([]db.CloudAccount, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestRepo) PatchGroupDisplayName(ctx context.Context, arg db.PatchGroupDisplayNameParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) PatchUser(ctx context.Context, arg db.PatchUserParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) SearchTag(ctx context.Context, arg db.SearchTagParams) ([]db.CloudAccount, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestRepo) UnAssignAccountFromOUs(ctx context.Context, cloudAccountID uuid.UUID) error {
	args := t.Called(ctx, cloudAccountID)

	return args.Error(0)
}

func (t *TestRepo) UpdateCloudAccount(ctx context.Context, arg db.UpdateCloudAccountParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) UpdateCloudAccountTagsDriftDetected(ctx context.Context, arg db.UpdateCloudAccountTagsDriftDetectedParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestRepo) WithTx(tx pgx.Tx) db.Repository {
	args := t.Called(tx)

	return args.Get(0).(db.Repository)
}

var _ db.Repository = (*TestRepo)(nil)
