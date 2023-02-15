package mocks

import (
	"context"

	"github.com/google/uuid"
	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/cli/fga"
)

type TestFGAAuthorizer struct {
	mock.Mock
}

func (t *TestFGAAuthorizer) WriteTuples(ctx context.Context, tuples []openfga.TupleKey) error {
	args := t.Called(ctx, tuples)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) WriteAssertions(ctx context.Context, authorizationModelID string, assertions []openfga.Assertion) error {
	args := t.Called(ctx, authorizationModelID, assertions)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) SetTypeDefinitions(ctx context.Context, authorizationModelID string) (string, error) {
	args := t.Called(ctx, authorizationModelID)

	return args.String(0), args.Error(1)
}

func (t *TestFGAAuthorizer) RunAssertions(ctx context.Context, typeDefinitionsContent string) (bool, error) {
	args := t.Called(ctx, typeDefinitionsContent)

	return args.Bool(0), args.Error(1)
}

func (t *TestFGAAuthorizer) ReplaceUsersInGroup(ctx context.Context, userIDs []uuid.UUID, groupID uuid.UUID) error {
	args := t.Called(ctx, userIDs, groupID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) Check(ctx context.Context, callerID uuid.UUID, relation fga.Relation, document fga.Document, objectID string) (bool, error) {
	args := t.Called(ctx, callerID, relation, document, objectID)

	return args.Bool(0), args.Error(1)
}

func (t *TestFGAAuthorizer) RemoveUser(ctx context.Context, userID uuid.UUID) error {
	args := t.Called(ctx, userID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) UserTuples(ctx context.Context, userID uuid.UUID, document string) ([]openfga.TupleKey, error) {
	args := t.Called(ctx, userID, document)

	return args.Get(0).([]openfga.TupleKey), args.Error(1)
}

func (t *TestFGAAuthorizer) CheckUserAlreadyExistsInOrganization(ctx context.Context, userID uuid.UUID) (bool, error) {
	args := t.Called(ctx, userID)

	return args.Bool(0), args.Error(1)
}

func (t *TestFGAAuthorizer) AddUserToOrganization(ctx context.Context, userID uuid.UUID) error {
	args := t.Called(ctx, userID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) RemoveUserFromOrganization(ctx context.Context, userID uuid.UUID) error {
	args := t.Called(ctx, userID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) CheckUserAlreadyExistsInGroup(ctx context.Context, userID, groupID uuid.UUID) (bool, error) {
	args := t.Called(ctx, userID, groupID)

	return args.Bool(0), args.Error(1)
}

func (t *TestFGAAuthorizer) AddUsersToGroup(ctx context.Context, userIDs []uuid.UUID, groupID uuid.UUID) error {
	args := t.Called(ctx, userIDs, groupID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) RemoveUserFromGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) error {
	args := t.Called(ctx, userID, groupID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) RemoveUsersInGroup(ctx context.Context, groupID uuid.UUID) error {
	args := t.Called(ctx, groupID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) CheckAccountAlreadyExistsInOrganization(ctx context.Context, accountID uuid.UUID) (bool, error) {
	args := t.Called(ctx, accountID)

	return args.Bool(0), args.Error(1)
}

func (t *TestFGAAuthorizer) AddAccountToOrganization(ctx context.Context, accountID uuid.UUID) error {
	args := t.Called(ctx, accountID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) CheckOrganizationalUnitRelationship(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) (bool, error) {
	args := t.Called(ctx, id, parentID)

	return args.Bool(0), args.Error(1)
}

func (t *TestFGAAuthorizer) AddOrganizationalUnit(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) error {
	args := t.Called(ctx, id, parentID)

	return args.Error(0)
}

func (t *TestFGAAuthorizer) RemoveOrganizationalUnitRelationships(ctx context.Context, id uuid.UUID, parentID uuid.NullUUID) error {
	args := t.Called(ctx, id, parentID)

	return args.Error(0)
}

var _ fga.Authorizer = (*TestFGAAuthorizer)(nil)
