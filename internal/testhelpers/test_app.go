package testhelpers

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/fga"
)

type TestApp struct {
	App        *application.App
	JS         *TestJS
	Repository *TestRepository
	FGAClient  *TestFGAAuthorizer
}

func NewTestApp() *TestApp {
	repository := new(TestRepository)
	js := new(TestJS)
	fgaClient := new(TestFGAAuthorizer)

	app := &application.App{
		Config:     application.Config{},
		JS:         js,
		Repository: repository,
		FGAClient:  fgaClient,
	}

	return &TestApp{
		App:        app,
		JS:         js,
		Repository: repository,
		FGAClient:  fgaClient,
	}
}

type TestFGAAuthorizer struct {
	mock.Mock
}

func (t *TestFGAAuthorizer) SetTypeDefinitions(ctx context.Context, typeDefinitionsContent string) (string, error) {
	args := t.Called(ctx, typeDefinitionsContent)

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

type TestRepository struct {
	mock.Mock
}

func (t *TestRepository) OrganizationalUnitsCloudAccounts(ctx context.Context, id []uuid.UUID) ([]db.CloudAccount, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestRepository) GetUserOrganizationalUnits(ctx context.Context, id uuid.UUID) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepository) GetAPIKeysOrganizationalUnits(ctx context.Context, id uuid.UUID) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepository) AssignCloudAccountToOrganizationalUnit(ctx context.Context, id, organizationalUnitID uuid.UUID) error {
	args := t.Called(ctx, id, organizationalUnitID)

	return args.Error(0)
}

func (t *TestRepository) UnAssignCloudAccountFromOrganizationalUnits(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepository) CreateTag(ctx context.Context, input db.CreateTagParams) (db.StandardTag, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.StandardTag), args.Error(1)
}

func (t *TestRepository) UpdateTag(ctx context.Context, input db.UpdateTagParams) (db.StandardTag, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.StandardTag), args.Error(1)
}

func (t *TestRepository) FindTag(ctx context.Context, id uuid.UUID) (db.StandardTag, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.StandardTag), args.Error(1)
}

func (t *TestRepository) GetTags(ctx context.Context) ([]db.StandardTag, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.StandardTag), args.Error(1)
}

func (t *TestRepository) CreateOrganizationalUnit(ctx context.Context, input db.CreateOrganizationalUnitParams) (db.OrganizationalUnit, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepository) GetOrganizationalUnits(ctx context.Context) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepository) FindOrganizationalUnit(ctx context.Context, id uuid.UUID) (db.OrganizationalUnit, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepository) GetOrganizationalUnitChildren(ctx context.Context, id uuid.UUID) ([]db.OrganizationalUnit, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.OrganizationalUnit), args.Error(1)
}

func (t *TestRepository) GetOrganizationalUnitCloudAccounts(ctx context.Context, id uuid.UUID) ([]db.CloudAccount, error) {
	args := t.Called(ctx, id)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestRepository) DeleteOrganizationalUnit(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepository) GetAPIKeys(ctx context.Context) ([]db.ApiKey, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.ApiKey), args.Error(1)
}

func (t *TestRepository) CreateAPIKey(ctx context.Context, input db.InsertAPIKeyParams) (db.ApiKey, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepository) FindAPIKey(ctx context.Context, id uuid.UUID) (db.ApiKey, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepository) GetAuditLogs(ctx context.Context) ([]db.AuditLog, []any, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.AuditLog), args.Get(1).([]any), args.Error(1)
}

func (t *TestRepository) GetAuditLogsForTarget(ctx context.Context, input db.GetAuditLogsForTargetParams) ([]db.AuditLog, []any, error) {
	args := t.Called(ctx, input)

	return args.Get(0).([]db.AuditLog), args.Get(1).([]any), args.Error(2)
}

func (t *TestRepository) InsertScimAPIKey(ctx context.Context, encodedHash string) (db.ApiKey, error) {
	args := t.Called(ctx, encodedHash)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepository) DeleteScimAPIKey(ctx context.Context) error {
	args := t.Called(ctx)

	return args.Error(0)
}

func (t *TestRepository) FindScimAPIKey(ctx context.Context) (db.ApiKey, error) {
	args := t.Called(ctx)

	return args.Get(0).(db.ApiKey), args.Error(1)
}

func (t *TestRepository) GetUsers(ctx context.Context, input db.GetUsersParams) ([]db.User, error) {
	args := t.Called(ctx, input)

	return args.Get(0).([]db.User), args.Error(1)
}

func (t *TestRepository) CreateAuditLog(ctx context.Context, input db.CreateAuditLogParams) (db.AuditLog, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.AuditLog), args.Error(1)
}

func (t *TestRepository) UpdateCloudAccountTagsDriftDetected(
	ctx context.Context,
	input db.UpdateCloudAccountTagsDriftDetectedParams,
) error {
	args := t.Called(ctx, input)

	return args.Error(0)
}

func (t *TestRepository) CreateOrInsertCloudAccount(
	ctx context.Context,
	input db.CreateOrInsertCloudAccountParams,
) (db.CloudAccount, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestRepository) CreateCloudTenant(ctx context.Context, input db.CreateCloudTenantParams) error {
	args := t.Called(ctx, input)

	return args.Error(0)
}

func (t *TestRepository) GetScimUsers(ctx context.Context, input db.GetScimUsersInput) (int64, []db.User, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(int64), args.Get(1).([]db.User), args.Error(2)
}

func (t *TestRepository) CreateUser(ctx context.Context, input db.CreateUserParams) (db.User, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepository) UpdateUser(ctx context.Context, id uuid.UUID, input db.UpdateUserParams) (db.User, error) {
	args := t.Called(ctx, id, input)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepository) ScimPatchUser(ctx context.Context, input db.PatchUserParams) error {
	args := t.Called(ctx, input)

	return args.Error(0)
}

func (t *TestRepository) Begin(ctx context.Context) (db.RepositoryQueries, error) {
	args := t.Called(ctx)

	return args.Get(0).(db.RepositoryQueries), args.Error(1)
}

func (t *TestRepository) Commit(ctx context.Context) error {
	args := t.Called(ctx)

	return args.Error(0)
}

func (t *TestRepository) Rollback(ctx context.Context) error {
	args := t.Called(ctx)

	return args.Error(0)
}

func (t *TestRepository) GetCloudTenants(ctx context.Context) ([]db.CloudTenant, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.CloudTenant), args.Error(1)
}

func (t *TestRepository) FindGroup(ctx context.Context, id string) (db.Group, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestRepository) CreateGroup(ctx context.Context, displayName string) (db.Group, error) {
	args := t.Called(ctx, displayName)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestRepository) DeleteGroup(ctx context.Context, id string) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestRepository) UpdateGroup(ctx context.Context, input db.PatchGroupDisplayNameParams) (db.Group, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestRepository) RemoveUserFromGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	args := t.Called(ctx, userID, groupID)

	return args.Error(0)
}

func (t *TestRepository) AddUserToGroup(ctx context.Context, userID, groupID uuid.UUID) error {
	args := t.Called(ctx, userID, groupID)

	return args.Error(0)
}

func (t *TestRepository) ReplaceUsersInGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error {
	args := t.Called(ctx, groupID, members)

	return args.Error(0)
}

func (t *TestRepository) AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []uuid.UUID) error {
	args := t.Called(ctx, groupID, members)

	return args.Error(0)
}

func (t *TestRepository) GetGroupMembership(ctx context.Context, idString string) ([]db.GetGroupMembershipRow, error) {
	args := t.Called(ctx, idString)

	return args.Get(0).([]db.GetGroupMembershipRow), args.Error(1)
}

func (t *TestRepository) GetGroups(ctx context.Context, params db.GetGroupsParams) (int64, []db.Group, error) {
	args := t.Called(ctx, params)

	return args.Get(0).(int64), args.Get(1).([]db.Group), args.Error(2)
}

func (t *TestRepository) FindUser(ctx context.Context, id string) (db.User, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepository) FindCloudAccount(
	ctx context.Context,
	input db.FindCloudAccountInput,
) (db.CloudAccount, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestRepository) UpdateCloudAccount(
	ctx context.Context,
	input db.UpdateCloudAccountParams,
) (db.CloudAccount, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestRepository) FindUserByUsername(ctx context.Context, username string) (db.User, error) {
	args := t.Called(ctx, username)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestRepository) SearchCloudAccounts(
	ctx context.Context,
	input db.SearchCloudAccountsInput,
) ([]db.CloudAccount, error) {
	args := t.Called(ctx, input)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

type TestJS struct {
	mock.Mock
}

func (t *TestJS) Publish(subj string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	args := t.Called(subj, data, opts)

	return args.Get(0).(*nats.PubAck), args.Error(1)
}

func (t *TestJS) PublishMsg(m *nats.Msg, opts ...nats.PubOpt) (*nats.PubAck, error) {
	args := t.Called(m, opts)

	return args.Get(0).(*nats.PubAck), args.Error(1)
}

func (t *TestJS) PublishAsync(subj string, data []byte, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	args := t.Called(subj, data, opts)

	return args.Get(0).(nats.PubAckFuture), args.Error(1)
}

func (t *TestJS) PublishMsgAsync(m *nats.Msg, opts ...nats.PubOpt) (nats.PubAckFuture, error) {
	args := t.Called(m, opts)

	return args.Get(0).(nats.PubAckFuture), args.Error(1)
}

func (t *TestJS) PublishAsyncPending() int {
	args := t.Called()

	return args.Get(0).(int)
}

func (t *TestJS) PublishAsyncComplete() <-chan struct{} {
	args := t.Called()

	return args.Get(0).(<-chan struct{})
}

func (t *TestJS) Subscribe(subj string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	args := t.Called(subj, cb, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) SubscribeSync(subj string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	args := t.Called(subj, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) ChanSubscribe(subj string, ch chan *nats.Msg, opts ...nats.SubOpt) (*nats.Subscription, error) {
	args := t.Called(subj, ch, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) ChanQueueSubscribe(subj,
	queue string,
	ch chan *nats.Msg,
	opts ...nats.SubOpt,
) (*nats.Subscription, error) {
	args := t.Called(subj, queue, ch, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) QueueSubscribe(subj, queue string,
	cb nats.MsgHandler,
	opts ...nats.SubOpt,
) (*nats.Subscription, error) {
	args := t.Called(subj, queue, cb, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) QueueSubscribeSync(subj, queue string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	args := t.Called(subj, queue, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) PullSubscribe(subj, durable string, opts ...nats.SubOpt) (*nats.Subscription, error) {
	args := t.Called(subj, durable, opts)

	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (t *TestJS) AddStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	args := t.Called(cfg, opts)

	return args.Get(0).(*nats.StreamInfo), args.Error(1)
}

func (t *TestJS) UpdateStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	args := t.Called(cfg, opts)

	return args.Get(0).(*nats.StreamInfo), args.Error(1)
}

func (t *TestJS) DeleteStream(name string, opts ...nats.JSOpt) error {
	args := t.Called(name, opts)

	return args.Error(0)
}

func (t *TestJS) StreamInfo(stream string, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	args := t.Called(stream, opts)

	return args.Get(0).(*nats.StreamInfo), args.Error(1)
}

func (t *TestJS) PurgeStream(name string, opts ...nats.JSOpt) error {
	args := t.Called(name, opts)

	return args.Error(0)
}

func (t *TestJS) StreamsInfo(opts ...nats.JSOpt) <-chan *nats.StreamInfo {
	args := t.Called(opts)

	return args.Get(0).(<-chan *nats.StreamInfo)
}

func (t *TestJS) StreamNames(opts ...nats.JSOpt) <-chan string {
	args := t.Called(opts)

	return args.Get(0).(<-chan string)
}

func (t *TestJS) GetMsg(name string, seq uint64, opts ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	args := t.Called(name, seq, opts)

	return args.Get(0).(*nats.RawStreamMsg), args.Error(1)
}

func (t *TestJS) DeleteMsg(name string, seq uint64, opts ...nats.JSOpt) error {
	args := t.Called(name, seq, opts)

	return args.Error(0)
}

func (t *TestJS) AddConsumer(stream string, cfg *nats.ConsumerConfig, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	args := t.Called(stream, cfg, opts)

	return args.Get(0).(*nats.ConsumerInfo), args.Error(1)
}

func (t *TestJS) UpdateConsumer(stream string,
	cfg *nats.ConsumerConfig,
	opts ...nats.JSOpt,
) (*nats.ConsumerInfo, error) {
	args := t.Called(stream, cfg, opts)

	return args.Get(0).(*nats.ConsumerInfo), args.Error(1)
}

func (t *TestJS) DeleteConsumer(stream, consumer string, opts ...nats.JSOpt) error {
	args := t.Called(stream, consumer, opts)

	return args.Error(0)
}

func (t *TestJS) ConsumerInfo(stream, name string, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	args := t.Called(stream, name, opts)

	return args.Get(0).(*nats.ConsumerInfo), args.Error(1)
}

func (t *TestJS) ConsumersInfo(stream string, opts ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	args := t.Called(stream, opts)

	return args.Get(0).(<-chan *nats.ConsumerInfo)
}

func (t *TestJS) ConsumerNames(stream string, opts ...nats.JSOpt) <-chan string {
	args := t.Called(stream, opts)

	return args.Get(0).(<-chan string)
}

func (t *TestJS) AccountInfo(opts ...nats.JSOpt) (*nats.AccountInfo, error) {
	args := t.Called(opts)

	return args.Get(0).(*nats.AccountInfo), args.Error(1)
}

func (t *TestJS) KeyValue(bucket string) (nats.KeyValue, error) {
	args := t.Called(bucket)

	return args.Get(0).(nats.KeyValue), args.Error(1)
}

func (t *TestJS) CreateKeyValue(cfg *nats.KeyValueConfig) (nats.KeyValue, error) {
	args := t.Called(cfg)

	return args.Get(0).(nats.KeyValue), args.Error(1)
}

func (t *TestJS) DeleteKeyValue(bucket string) error {
	args := t.Called(bucket)

	return args.Error(0)
}

func (t *TestJS) ObjectStore(bucket string) (nats.ObjectStore, error) {
	args := t.Called(bucket)

	return args.Get(0).(nats.ObjectStore), args.Error(1)
}

func (t *TestJS) CreateObjectStore(cfg *nats.ObjectStoreConfig) (nats.ObjectStore, error) {
	args := t.Called(cfg)

	return args.Get(0).(nats.ObjectStore), args.Error(1)
}

func (t *TestJS) DeleteObjectStore(bucket string) error {
	args := t.Called(bucket)

	return args.Error(0)
}
