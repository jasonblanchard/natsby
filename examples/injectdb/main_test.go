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
		if err != nil {
			panic(err)
		}

		_, err = nc.Request("store.greeting.set", []byte("hello"), time.Second*2)
		if err != nil {
			panic(err)
		}
		response, err := nc.Request("store.greeting.get", []byte(""), time.Second*2)

		So(err, ShouldBeNil)
		So(string(response.Data), ShouldEqual, "hello")
	})
}
