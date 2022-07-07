package helpers

import (
	"context"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
	"github.com/suse-skyscraper/skyscraper/internal/application"
	"github.com/suse-skyscraper/skyscraper/internal/db"
)

type TestAppResponse struct {
	App *application.App
	DB  *TestDB
	JS  *TestJS
}

func NewTestApp() *TestAppResponse {
	database := new(TestDB)
	js := new(TestJS)

	app := &application.App{
		Config: application.Config{},
		DB:     database,
		JS:     js,
	}

	return &TestAppResponse{
		App: app,
		DB:  database,
		JS:  js,
	}
}

type TestDB struct {
	mock.Mock
}

func (t *TestDB) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestDB) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := t.Called(ctx, id)

	return args.Error(0)
}

func (t *TestDB) GetGroup(ctx context.Context, id uuid.UUID) (db.Group, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestDB) CreateMembershipForUserAndGroup(
	ctx context.Context,
	arg db.CreateMembershipForUserAndGroupParams,
) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) DropMembershipForGroup(ctx context.Context, groupID uuid.UUID) error {
	args := t.Called(ctx, groupID)

	return args.Error(0)
}

func (t *TestDB) DropMembershipForUserAndGroup(ctx context.Context, arg db.DropMembershipForUserAndGroupParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) GetGroupMembership(ctx context.Context, groupID uuid.UUID) ([]db.GetGroupMembershipRow, error) {
	args := t.Called(ctx, groupID)

	return args.Get(0).([]db.GetGroupMembershipRow), args.Error(1)
}

func (t *TestDB) GetUser(ctx context.Context, id uuid.UUID) (db.User, error) {
	args := t.Called(ctx, id)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestDB) PatchGroupDisplayName(ctx context.Context, arg db.PatchGroupDisplayNameParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) CreateGroup(ctx context.Context, displayName string) (db.Group, error) {
	args := t.Called(ctx, displayName)

	return args.Get(0).(db.Group), args.Error(1)
}

func (t *TestDB) FindByUsername(ctx context.Context, username string) (db.User, error) {
	args := t.Called(ctx, username)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestDB) GetGroupCount(ctx context.Context) (int64, error) {
	args := t.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

func (t *TestDB) GetGroups(ctx context.Context, arg db.GetGroupsParams) ([]db.Group, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.Group), args.Error(1)
}

func (t *TestDB) PatchUser(ctx context.Context, arg db.PatchUserParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) UpdateUser(ctx context.Context, arg db.UpdateUserParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.User), args.Error(1)
}

func (t *TestDB) GetUserCount(ctx context.Context) (int64, error) {
	args := t.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

func (t *TestDB) GetUsers(ctx context.Context, arg db.GetUsersParams) ([]db.User, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.User), args.Error(1)
}

func (t *TestDB) CreateCloudTenant(ctx context.Context, arg db.CreateCloudTenantParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) CreateOrInsertCloudAccount(
	ctx context.Context,
	arg db.CreateOrInsertCloudAccountParams,
) (db.CloudAccount, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestDB) GetCloudAccount(ctx context.Context, arg db.GetCloudAccountParams) (db.CloudAccount, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudAccount), args.Error(1)
}

func (t *TestDB) GetCloudAllAccounts(ctx context.Context) ([]db.CloudAccount, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestDB) GetCloudAllAccountsForCloud(ctx context.Context, cloud string) ([]db.CloudAccount, error) {
	args := t.Called(ctx, cloud)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestDB) GetCloudAllAccountsForCloudAndTenant(
	ctx context.Context,
	arg db.GetCloudAllAccountsForCloudAndTenantParams,
) ([]db.CloudAccount, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).([]db.CloudAccount), args.Error(1)
}

func (t *TestDB) GetCloudTenant(ctx context.Context, arg db.GetCloudTenantParams) (db.CloudTenant, error) {
	args := t.Called(ctx, arg)

	return args.Get(0).(db.CloudTenant), args.Error(1)
}

func (t *TestDB) GetCloudTenants(ctx context.Context) ([]db.CloudTenant, error) {
	args := t.Called(ctx)

	return args.Get(0).([]db.CloudTenant), args.Error(1)
}

func (t *TestDB) UpdateCloudAccount(ctx context.Context, arg db.UpdateCloudAccountParams) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
}

func (t *TestDB) UpdateCloudAccountTagsDriftDetected(
	ctx context.Context,
	arg db.UpdateCloudAccountTagsDriftDetectedParams,
) error {
	args := t.Called(ctx, arg)

	return args.Error(0)
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
