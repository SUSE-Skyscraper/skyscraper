package helpers

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
	"github.com/suse-skyscraper/skyscraper/internal/scim/payloads"
)

type TestAppResponse struct {
	App        *application.App
	JS         *TestJS
	Repository *TestRepository
}

func NewTestApp() *TestAppResponse {
	repository := new(TestRepository)
	js := new(TestJS)

	app := &application.App{
		Config:     application.Config{},
		JS:         js,
		Repository: repository,
	}

	return &TestAppResponse{
		App:        app,
		JS:         js,
		Repository: repository,
	}
}

type TestRepository struct {
	mock.Mock
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

func (t *TestRepository) GetAuditLogs(ctx context.Context) ([]db.AuditLog, []db.User, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.AuditLog), args.Get(1).([]db.User), args.Error(2)
}

func (t *TestRepository) GetAuditLogsForTarget(
	ctx context.Context,
	input db.GetAuditLogsForTargetParams,
) ([]db.AuditLog, []db.User, error) {
	args := t.Called(ctx, input)

	return args.Get(0).([]db.AuditLog), args.Get(1).([]db.User), args.Error(2)
}

func (t *TestRepository) CreateAuditLog(ctx context.Context, input db.CreateAuditLogParams) (db.AuditLog, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.AuditLog), args.Error(1)
}

func (t *TestRepository) CreateTag(ctx context.Context, input db.CreateTagParams) (db.Tag, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.Tag), args.Error(1)
}

func (t *TestRepository) UpdateTag(ctx context.Context, input db.UpdateTagParams) (db.Tag, error) {
	args := t.Called(ctx, input)

	return args.Get(0).(db.Tag), args.Error(1)
}

func (t *TestRepository) FindTag(ctx context.Context, id uuid.UUID) (db.Tag, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.Tag), args.Error(1)
}

func (t *TestRepository) GetTags(ctx context.Context) ([]db.Tag, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.Tag), args.Error(1)
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

func (t *TestRepository) GetPolicies(ctx context.Context) ([]db.Policy, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.Policy), args.Error(1)
}

func (t *TestRepository) TruncatePolicies(ctx context.Context) error {
	args := t.Called(ctx)

	return args.Error(0)
}

func (t *TestRepository) CreatePolicy(ctx context.Context, input db.AddPolicyParams) error {
	args := t.Called(ctx, input)

	return args.Error(0)
}

func (t *TestRepository) RemovePolicy(ctx context.Context, input db.RemovePolicyParams) error {
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

func (t *TestRepository) ReplaceUsersInGroup(ctx context.Context,
	groupID uuid.UUID,
	members []payloads.MemberPatch,
) error {
	args := t.Called(ctx, groupID, members)

	return args.Error(0)
}

func (t *TestRepository) AddUsersToGroup(ctx context.Context, groupID uuid.UUID, members []payloads.MemberPatch) error {
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
