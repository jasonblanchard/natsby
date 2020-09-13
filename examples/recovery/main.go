package main

import (
	"fmt"

	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	engine, err := natsby.New(nc)
	if err != nil {
		panic(err)
	}

	logger := natsby.DefaultLogger()
	engine.Use(natsby.WithLogger(logger))
	engine.Use(natsby.WithCustomRecovery(func(c *natsby.Context, err interface{}) {
		logger.Error().Msg(fmt.Sprintf("%v", err))

		if c.Msg.Reply != "" {
			c.Conn.Publish(c.Msg.Reply, []byte("oops"))
		}
	}))
	// engine.Use(natsby.WithRecovery())

	engine.Subscribe("panic", natsby.WithByteReply(), func(c *natsby.Context) {
		panic("oops")
	})

	engine.Subscribe("ping", natsby.WithByteReply(), func(c *natsby.Context) {
		c.ByteReplyPayload = []byte("pong")
	})

	engine.Run(func() {
		logger.Info().Msg("Ready 🚀")
	})
}
