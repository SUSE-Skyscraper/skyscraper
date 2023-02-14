package helpers

import (
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type MockPubAckFuture struct {
	mock.Mock
}

func (m *MockPubAckFuture) Ok() <-chan *nats.PubAck {
	args := m.Called()

	return args.Get(0).(<-chan *nats.PubAck)
}

func (m *MockPubAckFuture) Err() <-chan error {
	args := m.Called()

	return args.Get(0).(<-chan error)
}

func (m *MockPubAckFuture) Msg() *nats.Msg {
	args := m.Called()

	return args.Get(0).(*nats.Msg)
}

var _ nats.PubAckFuture = (*MockPubAckFuture)(nil)
