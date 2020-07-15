package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

func main() {
	configureLogger := func(e *natsby.Engine) error {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zerolog.DurationFieldUnit = time.Second
		e.Logger = &logger
		return nil
	}

	configureNATS := func(e *natsby.Engine) error {
		nc, err := nats.Connect(nats.DefaultURL)
		c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
		if err != nil {
			return err
		}
		e.NatsConnection = nc
		e.NatsEncodedConnection = c
		return nil
	}

	engine, err := natsby.New(configureNATS, configureLogger)
	if err != nil {
		panic(err)
	}

	engine.Use(natsby.WithLogger())

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
