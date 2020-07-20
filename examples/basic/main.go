package main

import (
	"github.com/jasonblanchard/natsby"
)

func main() {
	engine, err := natsby.New()
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
