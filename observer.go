package natsby

import (
	"time"
)

// Observer interface for concrete observer types
type Observer interface {
	ObserveSubjectReceived(subject string)
	ObserveLatency(subject string, latency time.Duration)
	ObserveReply(subject string)
	ObserveError(subject string)
	Collect() error
}

// WithObserver Observes handler behavior and reports metrics
func WithObserver(o Observer) HandlerFunc {
	err := o.Collect()

	if err != nil {
		// TODO: Do something logical, here
		panic(err)
	}

	return func(c *Context) {
		o.ObserveSubjectReceived(c.Msg.Subject)

		latency := c.NextWithLatencyDuration()
		o.ObserveLatency(c.Msg.Subject, latency)

		_, ok := c.GetByteReplyPayload()

		if ok == true {
			o.ObserveReply(c.Msg.Subject)
		}

		// TODO: ObserveReply with JSON reply

		if c.Err != nil {
			o.ObserveError(c.Msg.Subject)
		}
	}
}
