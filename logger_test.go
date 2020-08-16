package natsby

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestDefaultLogger(t *testing.T) {
	logger := DefaultLogger()

	assert.IsType(t, zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}), *logger)
}

func TestWithLogger(t *testing.T) {
	var b bytes.Buffer
	logger := zerolog.New(&b)

	handler := WithLogger(&logger)

	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
			Reply:   "test.reply.inbox",
		},
	}

	handler(c)

	assert.Contains(t, b.String(), "test.subject")
	assert.Contains(t, b.String(), "test.reply.inbox")
	assert.Contains(t, b.String(), "latency")
}

func TestWithLoggerErr(t *testing.T) {
	var b bytes.Buffer
	logger := zerolog.New(&b)

	handler := WithLogger(&logger)

	c := &Context{
		Msg: &nats.Msg{
			Subject: "test.subject",
			Reply:   "test.reply.inbox",
		},
		Err: errors.New("Oops"),
	}

	handler(c)

	assert.Contains(t, b.String(), "test.subject")
	assert.Contains(t, b.String(), "test.reply.inbox")
	assert.Contains(t, b.String(), "latency")
	assert.Contains(t, b.String(), "Oops")
}
