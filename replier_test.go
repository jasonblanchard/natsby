package natsby

import (
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/nats-io/nats.go"
)

func TestWithByteReply(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	context := &Context{
		NatsConnection: nc,
		Msg: &nats.Msg{
			Reply: "reply.inbox",
		},
		ByteReplyPayload: []byte(""),
	}
	handler := WithByteReply()

	handler(context)

	assert.Equal(t, true, context.didReply)
}

func TestWithJsonReply(t *testing.T) {
	nc, _ := nats.Connect(nats.DefaultURL)
	encodedConnection, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	context := &Context{
		NatsConnection: nc,
		Msg: &nats.Msg{
			Reply: "reply.inbox",
		},
		ByteReplyPayload:      []byte(""),
		NatsEncodedConnection: encodedConnection,
	}
	handler := WithJSONReply()

	handler(context)

	assert.Equal(t, true, context.didReply)
}
