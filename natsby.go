package natsby

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
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
	Logger                *zerolog.Logger
	Subscribers           []*Subscriber
	middleware            HandlersChain
}

// New creates a new Router object
func New(options ...func(*Engine) error) (*Engine, error) {
	e := &Engine{}
	var err error

	for _, option := range options {
		err = option(e)
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

	e.Subscribers = append(e.Subscribers, s)
}

func (e *Engine) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(e.middleware) + len(handlers)
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, e.middleware)
	copy(mergedHandlers[len(e.middleware):], handlers)
	return mergedHandlers
}

// Run starts all the subscribers
func (e *Engine) Run() error {
	for _, subscriber := range e.Subscribers {
		handler := func(m *nats.Msg) {
			c := &Context{
				Msg:      m,
				handlers: subscriber.Handlers,
				engine:   e,
			}
			c.reset()
			c.Next()
		}

		e.NatsConnection.Subscribe(subscriber.Subject, handler)
	}

	return nil
}
