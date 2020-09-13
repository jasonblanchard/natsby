package natsby

// import (
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/nats-io/nats.go"
// 	"github.com/stretchr/testify/assert"
// )

// func TestNew(t *testing.T) {
// 	nc, _ := nats.Connect(nats.DefaultURL)
// 	engine, _ := New(nc)

// 	assert.Equal(t, os.Stdout, engine.OutWriter)
// 	assert.Equal(t, os.Stderr, engine.ErrWriter)
// 	assert.Equal(t, nc, engine.NatsConnection)
// }

// func TestNewOptions(t *testing.T) {
// 	nc, _ := nats.Connect(nats.DefaultURL)
// 	var c *nats.EncodedConn
// 	configureEncodedConnection := func(e *Engine) error {
// 		var err error
// 		c, err = nats.NewEncodedConn(e.NatsConnection, nats.JSON_ENCODER)
// 		if err != nil {
// 			return err
// 		}
// 		e.NatsEncodedConnection = c
// 		return nil
// 	}
// 	engine, _ := New(nc, configureEncodedConnection)

// 	assert.Equal(t, os.Stdout, engine.OutWriter)
// 	assert.Equal(t, os.Stderr, engine.ErrWriter)
// 	assert.Equal(t, nc, engine.NatsConnection)
// 	assert.Equal(t, c, engine.NatsEncodedConnection)
// }

// func TestUse(t *testing.T) {
// 	e := Engine{}
// 	handler := func(c *Context) {}

// 	e.Use(handler)

// 	assert.Len(t, e.middleware, 1)
// }

// func TestSubscribe(t *testing.T) {
// 	engine := Engine{}
// 	handler := func(c *Context) {}

// 	engine.Subscribe("test.subject", handler)

// 	assert.Equal(t, 1, len(engine.subscribers))
// }

// func TestRun(t *testing.T) {
// 	nc, _ := nats.Connect(nats.DefaultURL)
// 	engine, _ := New(nc)
// 	handler := func(c *Context) {
// 		c.ByteReplyPayload = []byte("pong")
// 	}
// 	engine.Subscribe("test.subject", WithByteReply(), handler)
// 	didCallCallback := false
// 	callback := func() {
// 		didCallCallback = true
// 	}

// 	go engine.Run(callback)

// 	// Let the listeners regiser
// 	time.Sleep(1 * time.Second)

// 	_, err := nc.Request("test.subject", []byte(""), 1*time.Second)
// 	if err != nil {
// 		panic(err)
// 	}

// 	engine.Shutdown()

// 	assert.True(t, didCallCallback)
// }

// func TestRunQueue(t *testing.T) {
// 	nc, _ := nats.Connect(nats.DefaultURL)
// 	configureQueueGroup := func(e *Engine) error {
// 		e.QueueGroup = "group"
// 		return nil
// 	}
// 	engine, _ := New(nc, configureQueueGroup)
// 	handler := func(c *Context) {}
// 	engine.Subscribe("test.subject", handler)

// 	go engine.Run()

// 	engine.Shutdown()

// 	assert.Equal(t, true, true, "Engine started and shutdown")
// }
