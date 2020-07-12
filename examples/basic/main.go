package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
		if err != nil {
			return err
		}
		e.NatsConnection = nc
		return nil
	}

	engine, err := natsby.New(configureNATS, configureLogger)
	if err != nil {
		panic(err)
	}

	engine.Use(natsby.WithLogger())

	engine.Subscribe("ping", natsby.WithByteReply(), func(c *natsby.Context) {
		c.ByteReplyPayload = []byte("pong")
	})

	err = engine.Run()

	if err != nil {
		panic(err)
	}

	fmt.Println("Ready 🚀")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go func() {
		// Wait for signal
		<-c
		engine.NatsConnection.Drain()
		os.Exit(0)
	}()
	runtime.Goexit()
}
