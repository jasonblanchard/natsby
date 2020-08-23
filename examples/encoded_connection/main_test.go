package main

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		panic(err)
	}

	type pinger struct {
		Greeting string
	}

	ping := &pinger{
		Greeting: "ping",
	}

	pong := &pinger{}

	err = c.Request("ping", ping, pong, time.Second*2)

	assert.Nil(t, err)
	assert.Equal(t, "pong", pong.Greeting)
}
