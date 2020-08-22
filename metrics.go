package natsby

import (
	"time"
)

// Collector metrics collector
type Collector interface {
	CollectSubjectReceived(subject string)
	CollectLatency(subject string, latency time.Duration)
	CollectReply(subject string)
	CollectError(subject string)
	Collect() error
}

// Observer describes when collectors should collect data
type Observer interface {
	ObserveSubjectReceived(c *Context)
	ObserveLatency(c *Context, latency time.Duration)
	ObserveReply(c *Context)
	ObserveError(c *Context)
	Observe() error
}

// WithMetrics Observes handler behavior and reports metrics
func WithMetrics(o Observer) HandlerFunc {
	err := o.Observe()

	if err != nil {
		// TODO: Do something logical, here
		panic(err)
	}

	return func(c *Context) {
		o.ObserveSubjectReceived(c)
		latency := c.NextWithLatencyDuration()
		o.ObserveLatency(c, latency)
		o.ObserveReply(c)
		o.ObserveError(c)
	}
}
