package natsby

import (
	"github.com/nats-io/nats.go"
)

// HandlerFunc defines handler used by middleware as return value
type HandlerFunc func(*Context)

// HandlersChain HandlerFunc array
type HandlersChain []HandlerFunc

// Subscriber respresents a subscriber to be set up in Run()
type Subscriber struct {
	Subject  string
	Handlers HandlersChain
}

// Engine framework instance
type Engine struct {
	NatsConnection        *nats.Conn
	NatsEncodedConnection *nats.EncodedConn
	subscribers           []*Subscriber
	middleware            HandlersChain
	done                  chan bool
	queueGroup            string // TODO: Make this configurable
}

// New creates a new Router object
// TODO: Make connection be a required first argument
func New(options ...func(*Engine) error) (*Engine, error) {
	e := &Engine{
		done: make(chan bool),
	}
	var err error

	for _, option := range options {
		err = option(e)
	}

	if e.NatsConnection == nil {
		nc, err := nats.Connect(nats.DefaultURL)
		if err != nil {
			return e, err
		}
		e.NatsConnection = nc
	}

	return e, err
}

// Use adds global middleware to the engine which will be called for every subscription
func (e *Engine) Use(middleware ...HandlerFunc) {
	e.middleware = append(e.middleware, middleware...)
}

// Subscribe adds a subscriber to the NATS instance with middleware
func (e *Engine) Subscribe(subject string, handlers ...HandlerFunc) {
	s := &Subscriber{
		Subject:  subject,
		Handlers: e.combineHandlers(handlers),
	}

	e.subscribers = append(e.subscribers, s)
}

func (e *Engine) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(e.middleware) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, e.middleware)
	copy(mergedHandlers[len(e.middleware):], handlers)
	return mergedHandlers
}

// Run starts all the subscribers and blocks
func (e *Engine) Run(cbs ...func()) error {
	for _, subscriber := range e.subscribers {
		func(subscriber *Subscriber) {
			handler := func(m *nats.Msg) {
				c := &Context{
					Msg:      m,
					handlers: subscriber.Handlers,
					Engine:   e,
					Keys:     make(map[string]interface{}),
				}
				c.reset()
				c.Next()
			}

			if e.queueGroup == "" {
				e.NatsConnection.Subscribe(subscriber.Subject, handler)
				return
			}
			e.NatsConnection.QueueSubscribe(subscriber.Subject, e.queueGroup, handler)
		}(subscriber)
	}

	for _, cb := range cbs {
		cb()
	}

	<-e.done

	e.NatsConnection.Drain()

	return nil
}

// Shutdown terminates all listeners and drains connections
func (e *Engine) Shutdown() {
	e.done <- true
}
