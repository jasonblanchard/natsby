package natsby

import (
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	engine, _ := New(nc)

	assert.Equal(t, os.Stdout, engine.OutWriter)
	assert.Equal(t, os.Stderr, engine.ErrWriter)
	assert.Equal(t, nc, engine.NatsConnection)
}

func TestNewOptions(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	var c *nats.EncodedConn
	configureEncodedConnection := func(e *Engine) error {
		var err error
		c, err = nats.NewEncodedConn(e.NatsConnection, nats.JSON_ENCODER)
		if err != nil {
			return err
		}
		e.NatsEncodedConnection = c
		return nil
	}
	engine, _ := New(nc, configureEncodedConnection)

	assert.Equal(t, os.Stdout, engine.OutWriter)
	assert.Equal(t, os.Stderr, engine.ErrWriter)
	assert.Equal(t, nc, engine.NatsConnection)
	assert.Equal(t, c, engine.NatsEncodedConnection)
}

func TestSubscribe(t *testing.T) {
	engine := Engine{}
	handler := func(c *Context) {}

	engine.Subscribe("test.subject", handler)

	assert.Equal(t, 1, len(engine.subscribers))
}
