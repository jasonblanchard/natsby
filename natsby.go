package natsby

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// Event event sent to a controller handler
type Event struct {
	*nats.Msg
}

// EventResult data about how the event was handled
type EventResult struct{}

// Controller responds to message events
type Controller interface {
	Handle(e *Event) (EventResult, error) // TODO: Pass this something
	// Setup()
}

// Subscriber controller and subject it should respond to
type Subscriber struct {
	Subject string
	Ctrl    Controller
}

// Engine root framework instance
type Engine struct {
	*nats.Conn
	Subscribers []*Subscriber
	done        chan bool
	Logger      *zerolog.Logger
}

// Subscribe map a controller to a subject
func (e *Engine) Subscribe(subject string, ctrl Controller) {
	s := &Subscriber{
		Subject: subject,
		Ctrl:    ctrl,
	}
	e.Subscribers = append(e.Subscribers, s)
}

// Run starts all subscriber controllers and blocks
func (e *Engine) Run(callbacks ...func()) error {
	for _, subscriber := range e.Subscribers {
		func(subscriber *Subscriber) {
			handlerFunc := func(m *nats.Msg) {
				e := &Event{
					Msg: m,
				}

				subscriber.Ctrl.Handle(e)
			}

			e.Conn.Subscribe(subscriber.Subject, handlerFunc)
		}(subscriber)
	}

	for _, cb := range callbacks {
		cb()
	}

	<-e.done

	e.Conn.Drain()

	return nil
}
