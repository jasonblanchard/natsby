package main

import (
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

	engine.Subscribe("ping", natsby.WithByteReply(), func(c *natsby.Context) {
		c.ByteReplyPayload = []byte("pong")
	})

	engine.Run(func() {
		logger.Info().Msg("Ready 🚀")
	})

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGINT)
	// go func() {
	// 	// Wait for signal
	// 	<-c
	// 	engine.NatsConnection.Drain()
	// 	os.Exit(0)
	// }()
	// runtime.Goexit()
}
