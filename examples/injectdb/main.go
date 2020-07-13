package main

import (
	"errors"
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

type dB struct {
	state map[string]string
}

func newDB() *dB {
	db := &dB{
		state: make(map[string]string, 0),
	}

	return db
}

func (db *dB) Set(key, value string) {
	db.state[key] = value
}

func (db *dB) Get(key string) string {
	return db.state[key]
}

func withDb(db *dB) natsby.HandlerFunc {
	return func(c *natsby.Context) {
		c.Set("db", db)
		c.Next()
	}
}

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

	db := newDB()

	engine.Subscribe("store.greeting.set", natsby.WithByteReply(), withDb(db), func(c *natsby.Context) {
		db, ok := c.Get("db").(*dB)
		if ok != true {
			c.Err = errors.New("DB not what I expected")
			return
		}

		db.Set("greeting", string(c.Msg.Data))
	})

	engine.Subscribe("store.greeting.get", natsby.WithByteReply(), withDb(db), func(c *natsby.Context) {
		db, ok := c.Get("db").(*dB)
		if ok != true {
			c.Err = errors.New("DB not what I expected")
			return
		}
		greeting := db.Get("greeting")
		c.ByteReplyPayload = []byte(greeting)
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
