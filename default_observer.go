package natsby

import "time"

// DefaultObserver Default concrete Observer implementation
type DefaultObserver struct {
	collector Collector
}

// NewDefaultObserver creates a default observer instance
func NewDefaultObserver(c Collector) *DefaultObserver {
	o := &DefaultObserver{
		collector: c,
	}

	return o
}

// ObserveSubjectReceived collect metrics when subject is received
func (o *DefaultObserver) ObserveSubjectReceived(c *Context) {
	o.collector.CollectSubjectReceived(c.Msg.Subject)
}

// ObserveLatency collect metrics on handler chain latency
func (o *DefaultObserver) ObserveLatency(c *Context, latency time.Duration) {
	o.collector.CollectLatency(c.Msg.Subject, latency)
}

// ObserveReply collect metrics on replies
func (o *DefaultObserver) ObserveReply(c *Context) {
	_, ok := c.GetByteReplyPayload()
	if ok == true {
		o.collector.CollectReply(c.Msg.Subject)
	}

	// TODO: ObserveReply with JSON reply
}

// ObserveError collect metrics on replies
func (o *DefaultObserver) ObserveError(c *Context) {
	if c.Err != nil {
		o.collector.CollectError(c.Msg.Subject)
	}
}

// Observe start observing and collecting metrics
func (o *DefaultObserver) Observe() error {
	o.collector.Collect()
	return nil
}
