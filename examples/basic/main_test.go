package main

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	response, err := nc.Request("ping", []byte(""), time.Second*2)

	assert.Nil(t, err)
	assert.Equal(t, "pong", string(response.Data))
}
