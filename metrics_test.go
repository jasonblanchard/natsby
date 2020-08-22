package natsby

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type TestMetricsMockObserver struct {
	mock.Mock
}

func (m *TestMetricsMockObserver) ObserveSubjectReceived(c *Context) {
	m.Called(c)
}

func (m *TestMetricsMockObserver) ObserveLatency(c *Context, latency time.Duration) {
	m.Called(c, latency)
}

func (m *TestMetricsMockObserver) ObserveReply(c *Context) {
	m.Called(c)
}

func (m *TestMetricsMockObserver) ObserveError(c *Context) {
	m.Called(c)
}

func (m *TestMetricsMockObserver) Observe() error {
	return nil
}

func TestWithMetrics(t *testing.T) {
	observer := new(TestMetricsMockObserver)
	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
		},
	}

	observer.On("ObserveSubjectReceived", c).Return(true, nil)
	observer.On("ObserveLatency", c, mock.Anything).Return(true, nil)
	observer.On("ObserveReply", c).Return(true, nil)
	observer.On("ObserveError", c).Return(true, nil)
	observer.On("Observe").Return(true, nil)

	handler := WithMetrics(observer)

	handler(c)
}
