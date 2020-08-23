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

	_, err = nc.Request("store.greeting.set", []byte("hello"), time.Second*2)
	if err != nil {
		panic(err)
	}
	response, err := nc.Request("store.greeting.get", []byte(""), time.Second*2)

	assert.Nil(t, err)
	assert.Equal(t, "hello", string(response.Data))
}
