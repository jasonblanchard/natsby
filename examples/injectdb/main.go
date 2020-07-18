package main

import (
	"errors"
	"fmt"

	"github.com/jasonblanchard/natsby"
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
	engine, err := natsby.New()
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

	engine.Run(func() {
		fmt.Println("Ready 🚀")
	})
}
