package natsby

import (
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type MockObserver struct {
	mock.Mock
}

func (m *MockObserver) ObserveSubjectReceived(subject string) {
	m.Called(subject)
}

func (m *MockObserver) ObserveLatency(subject string, latency time.Duration) {
	m.Called(subject, latency)
}

func (m *MockObserver) ObserveReply(subject string) {
	m.Called(subject)
}

func (m *MockObserver) ObserveError(subject string) {
	m.Called(subject)
}

func (m *MockObserver) Collect() error {
	// m.Called()
	return nil
}

func TestWithObserver(t *testing.T) {
	observer := new(MockObserver)
	handler := WithObserver(observer)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
	}

	observer.On("ObserveSubjectReceived", "test.subject").Return(true, nil)
	observer.On("ObserveLatency", mock.Anything, mock.Anything).Return(true, nil)

	handler(c)

	observer.AssertExpectations(t)
}

func TestWithObserverWithByteReply(t *testing.T) {
	observer := new(MockObserver)
	handler := WithObserver(observer)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
		didReply:         true,
		ByteReplyPayload: []byte(""),
	}

	observer.On("ObserveSubjectReceived", "test.subject").Return(true, nil)
	observer.On("ObserveLatency", mock.Anything, mock.Anything).Return(true, nil)
	observer.On("ObserveReply", "test.subject").Return(true, nil)

	handler(c)

	observer.AssertExpectations(t)
}

func TestWithError(t *testing.T) {
	observer := new(MockObserver)
	handler := WithObserver(observer)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
		Err: errors.New("oops"),
	}

	observer.On("ObserveSubjectReceived", "test.subject").Return(true, nil)
	observer.On("ObserveLatency", mock.Anything, mock.Anything).Return(true, nil)
	observer.On("ObserveError", "test.subject").Return(true, nil)

	handler(c)

	observer.AssertExpectations(t)
}
