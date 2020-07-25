package main

import (
	"fmt"

	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	configureEncodedConnection := func(e *natsby.Engine) error {
		c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			return err
		}
		e.NatsEncodedConnection = c
		return nil
	}

	engine, err := natsby.New(nc, configureEncodedConnection)
	if err != nil {
		panic(err)
	}

	engine.Use(natsby.WithLogger(natsby.DefaultLogger()))

	engine.Subscribe("ping", natsby.WithJSONReply(), func(c *natsby.Context) {
		type pinger struct {
			Greeting string
		}

		payload := &pinger{
			Greeting: "pong",
		}

		c.JSONReplyPayload = payload
	})

	engine.Run(func() {
		fmt.Println("Ready 🚀")
	})
}
