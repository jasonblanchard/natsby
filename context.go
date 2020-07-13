package natsby

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// Context context that's passed through handlers and middleware
type Context struct {
	Msg              *nats.Msg
	handlers         HandlersChain
	ByteReplyPayload []byte
	JSONReplyPayload interface{}
	index            int8
	engine           *Engine
	Err              error
	Logger           *zerolog.Logger
	Keys             map[string]interface{}
}

// Next to be called in middleware to invoke the middleware chain
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) reset() {
	c.index = -1
}

// Set sets arbitrary value that will be available in the context map
func (c *Context) Set(k string, v interface{}) {
	c.Keys[k] = v
}

// Get gets arbirary value from the context map
func (c *Context) Get(k string) interface{} {
	return c.Keys[k]
}
