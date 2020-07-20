package main

import (
	"fmt"

	"github.com/jasonblanchard/natsby"
)

func main() {
	engine, err := natsby.New()
	if err != nil {
		panic(err)
	}

	engine.Use(natsby.WithLogger())
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
