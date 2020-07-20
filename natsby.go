package natsby

import (
	"os"
	"time"

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
	done                  chan bool
	queueGroup            string
}

// New creates a new Router object
func New(options ...func(*Engine) error) (*Engine, error) {
	e := &Engine{
		done: make(chan bool),
	}
	var err error

	for _, option := range options {
		err = option(e)
	}

	// TODO: Move all logger initialization to middleware?
	if e.Logger == nil {
		logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		zerolog.DurationFieldUnit = time.Second
		e.Logger = &logger
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

	e.Subscribers = append(e.Subscribers, s)
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
	for _, subscriber := range e.Subscribers {
		func(subscriber *Subscriber) {
			handler := func(m *nats.Msg) {
				c := &Context{
					Msg:      m,
					handlers: subscriber.Handlers,
					Engine:   e,
					Logger:   e.Logger,
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
	e.Logger.Info().Msg("Closing natsby listeners")
	e.done <- true
}
