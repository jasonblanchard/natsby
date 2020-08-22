package natsby

import (
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type TestObserverMockCollector struct {
	mock.Mock
}

func (m *TestObserverMockCollector) CollectSubjectReceived(subject string) {
	m.Called(subject)
}

func (m *TestObserverMockCollector) CollectLatency(subject string, latency time.Duration) {
	m.Called(subject, latency)
}

func (m *TestObserverMockCollector) CollectReply(subject string) {
	m.Called(subject)
}

func (m *TestObserverMockCollector) CollectError(subject string) {
	m.Called(subject)
}

func (m *TestObserverMockCollector) Collect() error {
	return nil
}

func TestObserveSubjectReceived(t *testing.T) {
	collector := new(TestObserverMockCollector)
	observer := NewDefaultObserver(collector)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
	}

	collector.On("CollectSubjectReceived", "test.subject").Return(true, nil)

	observer.ObserveSubjectReceived(c)

	collector.AssertExpectations(t)
}

func TestObserveLatency(t *testing.T) {
	collector := new(TestObserverMockCollector)
	observer := NewDefaultObserver(collector)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
	}

	collector.On("CollectLatency", "test.subject", mock.Anything).Return(true, nil)

	observer.ObserveLatency(c, 2*time.Second)

	collector.AssertExpectations(t)
}

func TestObserveByteReply(t *testing.T) {
	collector := new(TestObserverMockCollector)
	observer := NewDefaultObserver(collector)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
		didReply:         true,
		ByteReplyPayload: []byte("reply"),
	}

	collector.On("CollectReply", "test.subject").Return(true, nil)

	observer.ObserveReply(c)

	collector.AssertExpectations(t)
}

func TestObserveError(t *testing.T) {
	collector := new(TestObserverMockCollector)
	observer := NewDefaultObserver(collector)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
		Err: errors.New("oops"),
	}

	collector.On("CollectError", "test.subject").Return(true, nil)

	observer.ObserveError(c)

	collector.AssertExpectations(t)
}

func TestObserve(t *testing.T) {
	collector := new(TestObserverMockCollector)
	observer := NewDefaultObserver(collector)

	collector.On("Collect").Return(true, nil)

	observer.Observe()
}
