package natsby

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func panicHandler(c *Context) {
	panic("oops")
}

func TestWithRecovery(t *testing.T) {
	context := &Context{
		handlers:  HandlersChain{panicHandler},
		errWriter: bytes.NewBufferString(""),
	}

	context.reset()

	handler := WithRecovery()

	handler(context)

	out, _ := ioutil.ReadAll(context.errWriter)
	assert.Contains(t, string(out), "panic recovered oops")
}

func TestWithCusomRecovery(t *testing.T) {
	context := &Context{
		handlers:  HandlersChain{panicHandler},
		errWriter: bytes.NewBufferString(""),
	}

	context.reset()

	var customMessage string

	handler := WithCustomRecovery(func(c *Context, err interface{}) {
		customMessage = "Set in custom recovery fn"
	})

	handler(context)

	out, _ := ioutil.ReadAll(context.errWriter)
	assert.Contains(t, string(out), "panic recovered oops")
	assert.Equal(t, customMessage, "Set in custom recovery fn")
}
