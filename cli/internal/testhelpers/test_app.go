package testhelpers

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"

	"github.com/suse-skyscraper/skyscraper/cli/application"
	"github.com/suse-skyscraper/skyscraper/cli/internal/db"
	"github.com/suse-skyscraper/skyscraper/cli/internal/fga"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	openfga "github.com/openfga/go-sdk"
	"github.com/stretchr/testify/mock"
)

type TestApp struct {
	App          *application.App
	JS           *TestJS
	Repo         *TestRepo
	Searcher     *TestSearcher
	FGAClient    *TestFGAAuthorizer
	PostgresPool pgxmock.PgxPoolIface
}

func (t *TestApp) Close() {
	t.PostgresPool.Close()
}

func NewTestApp() (*TestApp, error) {
	js := new(TestJS)
	fgaClient := new(TestFGAAuthorizer)
	repo := new(TestRepo)
	searcher := new(TestSearcher)
	pool, err := pgxmock.NewPool()
	if err != nil {
		return nil, err
	}

	app := &application.App{
		Config:       application.Config{},
		JS:           js,
		FGAClient:    fgaClient,
		Repo:         repo,
		PostgresPool: pool,
		Searcher:     searcher,
	}

	return &TestApp{
		App:          app,
		JS:           js,
		FGAClient:    fgaClient,
		PostgresPool: pool,
		Repo:         repo,
		Searcher:     searcher,
	}, nil
}

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

type TestJS struct {
	mock.Mock
}

func (t *TestJS) StreamNameBySubject(s string, opt ...nats.JSOpt) (string, error) {
	args := t.Called(s, opt)

	return args.String(0), args.Error(1)
}

func (t *TestJS) KeyValueStoreNames() <-chan string {
	args := t.Called()

	return args.Get(0).(<-chan string)
}

func (t *TestJS) KeyValueStores() <-chan nats.KeyValueStatus {
	args := t.Called()

	return args.Get(0).(<-chan nats.KeyValueStatus)
}

func (t *TestJS) ObjectStores(opts ...nats.ObjectOpt) <-chan nats.ObjectStoreStatus {
	args := t.Called(opts)

	return args.Get(0).(<-chan nats.ObjectStoreStatus)
}

func (t *TestJS) Streams(opts ...nats.JSOpt) <-chan *nats.StreamInfo {
	args := t.Called(opts)

	return args.Get(0).(<-chan *nats.StreamInfo)
}

func (t *TestJS) GetLastMsg(name, subject string, opts ...nats.JSOpt) (*nats.RawStreamMsg, error) {
	args := t.Called(name, subject, opts)

	return args.Get(0).(*nats.RawStreamMsg), args.Error(1)
}

func (t *TestJS) SecureDeleteMsg(name string, seq uint64, opts ...nats.JSOpt) error {
	args := t.Called(name, seq, opts)

	return args.Error(0)
}

func (t *TestJS) Consumers(stream string, opts ...nats.JSOpt) <-chan *nats.ConsumerInfo {
	args := t.Called(stream, opts)

	return args.Get(0).(<-chan *nats.ConsumerInfo)
}

func (t *TestJS) ObjectStoreNames(opts ...nats.ObjectOpt) <-chan string {
	args := t.Called(opts)

	return args.Get(0).(<-chan string)
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

type TestSearcher struct {
	mock.Mock
}

func (t *TestSearcher) SearchCloudAccounts(ctx context.Context, input db.SearchCloudAccountsInput) ([]db.CloudAccount, error) {
	args := t.Called(ctx, input)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}
