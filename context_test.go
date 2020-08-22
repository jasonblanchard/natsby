package natsby

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	setOutHandlerCalled := false

	setOutHandler := func(c *Context) {
		setOutHandlerCalled = true
	}
	c := Context{
		handlers: HandlersChain{setOutHandler},
	}
	c.reset()

	c.Next()
	assert.Equal(t, true, setOutHandlerCalled)
}

func TestNextWithLatencyDuration(t *testing.T) {
	setOutHandlerCalled := false

	setOutHandler := func(c *Context) {
		setOutHandlerCalled = true
	}
	c := Context{
		handlers: HandlersChain{setOutHandler},
	}
	c.reset()

	latency := c.NextWithLatencyDuration()
	assert.Equal(t, true, setOutHandlerCalled)
	var d time.Duration
	assert.IsType(t, d, latency)
}

func TestKeys(t *testing.T) {
	c := Context{}

	c.Set("test", "hello")
	out := c.Get("test").(string)
	assert.Equal(t, "hello", out)
}

func TestGetByteReplyPayload(t *testing.T) {
	c := &Context{}
	var payload []byte

	payload, ok := c.GetByteReplyPayload()

	assert.Equal(t, ok, false)
	assert.Equal(t, []byte(""), payload)

	c.didReply = true
	c.ByteReplyPayload = []byte("reply")

	payload, ok = c.GetByteReplyPayload()
	assert.Equal(t, ok, true)
	assert.Equal(t, []byte("reply"), payload)
}
