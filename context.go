package natsby

import (
	"io"
	"time"

	"github.com/nats-io/nats.go"
)

// Context context that's passed through handlers and middleware
type Context struct {
	Msg                   *nats.Msg
	handlers              HandlersChain
	ByteReplyPayload      []byte
	JSONReplyPayload      interface{}
	didReply              bool
	index                 int8
	NatsConnection        *nats.Conn
	NatsEncodedConnection *nats.EncodedConn
	Err                   error
	Keys                  map[string]interface{}
	outWriter             io.ReadWriter
	errWriter             io.ReadWriter
}

// Next to be called in middleware to invoke the middleware chain
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// NextWithLatencyDuration Calls Next() and returns duration of execution
func (c *Context) NextWithLatencyDuration() time.Duration {
	start := time.Now()
	c.Next()
	end := time.Now()
	return end.Sub(start)
}

func (c *Context) reset() {
	c.index = -1
}

// Set sets arbitrary value that will be available in the context map
func (c *Context) Set(k string, v interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[k] = v
}

// Get gets arbirary value from the context map
func (c *Context) Get(k string) interface{} {
	return c.Keys[k]
}

// GetByteReplyPayload getter for byte reply payload with metadata about if it was set
func (c *Context) GetByteReplyPayload() ([]byte, bool) {
	if c.didReply == false {
		return []byte(""), false
	}
	return c.ByteReplyPayload, true
}

// TODO: GetJSONReplyPayload
