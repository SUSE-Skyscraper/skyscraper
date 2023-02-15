package mocks

import (
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

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

var _ nats.JetStreamContext = (*TestJS)(nil)
