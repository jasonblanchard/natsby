package natsby

import (
	"io"
	"os"

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
	QueueGroup            string
	OutWriter             io.ReadWriter
	ErrWriter             io.ReadWriter
}

// New creates a new Router object
func New(nc *nats.Conn, options ...func(*Engine) error) (*Engine, error) {
	e := &Engine{
		done: make(chan bool),
	}
	var err error

	e.OutWriter = os.Stdout
	e.ErrWriter = os.Stderr
	e.NatsConnection = nc

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
func (e *Engine) Run(callbacks ...func()) error {
	for _, subscriber := range e.subscribers {
		func(subscriber *Subscriber) {
			handler := func(m *nats.Msg) {
				c := &Context{
					Msg:                   m,
					handlers:              subscriber.Handlers,
					NatsConnection:        e.NatsConnection,
					NatsEncodedConnection: e.NatsEncodedConnection,
					Keys:                  make(map[string]interface{}),
					outWriter:             e.OutWriter,
					errWriter:             e.ErrWriter,
				}
				c.reset()
				c.Next()
			}

			if e.QueueGroup == "" {
				e.NatsConnection.Subscribe(subscriber.Subject, handler)
				return
			}
			e.NatsConnection.QueueSubscribe(subscriber.Subject, e.QueueGroup, handler)
		}(subscriber)
	}

	for _, cb := range callbacks {
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
