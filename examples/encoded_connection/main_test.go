package main

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {
	Convey("it works", t, func() {
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

		So(err, ShouldBeNil)
		So(pong.Greeting, ShouldEqual, "pong")
	})
}
