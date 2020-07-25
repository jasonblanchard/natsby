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

	engine.Use(natsby.WithLogger(natsby.DefaultLogger()))
	engine.Use(natsby.WithPrometheus(&natsby.WithPrometheusInput{
		Port: "2112",
	}))

	engine.Subscribe("ping", natsby.WithByteReply(), func(c *natsby.Context) {
		// time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		c.ByteReplyPayload = []byte("pong")
	})

	engine.Run(func() {
		fmt.Println("Ready 🚀")
	})
}
