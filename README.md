# Natsby
![ci status](https://github.com/jasonblanchard/natsby/workflows/CI/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/jasonblanchard/natsby/badge.svg?branch=master)](https://coveralls.io/github/jasonblanchard/natsby?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/jasonblanchard/natsby)](https://goreportcard.com/report/github.com/jasonblanchard/natsby)


Natsby enables NATS-driven services.

## Installation
Install the dependency:
```bash
go get -u github.com/jasonblanchard/natsby
```

Import in your code:
```go
import "github.com/nats-io/nats.go"
```

## Quickstart
Assuming a NATS server running at `http://localhost:4222`, initialize the NATS connection and pass it to Natsby:

```go
package main

import (
	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	engine, err := natsby.New(nc)
	if err != nil {
		panic(err)
	}

	logger := natsby.DefaultLogger()
	engine.Use(natsby.WithLogger(logger)) // (1)

	engine.Subscribe("ping", natsby.WithByteReply(), func(c *natsby.Context) { // (2)
		c.ByteReplyPayload = []byte("pong")
	})

	engine.Run(func() { // (3)
		logger.Info().Msg("Ready 🚀")
	})
}
```

This will:
1. Set up the default logger as global middleware. This will be run on each received message.
2. Register a middleware chain to be called when messages are received on the "ping" subject. The handlers will be called right-to-left. `WithByteReply` is a built-in middleware that will publish the payload in `ByteReplyPayload` to the message reply subject.
3. Start all the subscriptions. By default, this is blocking.

See [the examples directory](./examples) for more sample code.

## Middleware
Middleware is a chain of functions that are run when a message is received on a subscribed topic.

Middleware can be added globally using `engine.Use()` (i.e. it is called on every subscription) or per subscription as functions passed to `engine.Subscribe()`.

Middleware functions are passed a [`Context`](./context.go) instance unique to that subscription handler invocation. The `Context` instance can be mutated to provide data for subsequent middleware via specific setters on the `Context` struct or generic key/value pairs via `Context.Get()`/`Context.Set()`.

### Built-in Middleware
#### Logging
Logs information as messages are received and handled. This middleware uses [zerolog](https://github.com/rs/zerolog) (this may change in the future).

You can bring your own zerolog instance or use the default logger:

```go
logger := natsby.DefaultLogger()
engine.Use(natsby.WithLogger(logger))
```

#### Replier
Publishes the payload in `context.ByteReplyPayload` or `context.JSONReplyPayload` to the `context.Msg.Reply` subject if it exists.

#### Recovery
Catches `panic()`s and converts it to an error so that the process does not crash. By default, it logs the error and stack trace:
```go
engine.Use(natsby.WithRecovery())
```

You can also bring your own recovery handler:
```go
	engine.Use(natsby.WithCustomRecovery(func(c *natsby.Context, err interface{}) {
		logger.Error().Msg(fmt.Sprintf("%v", err))

		if c.Msg.Reply != "" {
			c.NatsConnection.Publish(c.Msg.Reply, []byte("oops"))
		}
  }))
```

#### Metrics
Collect observability metrics for various aspects of Natsby. The "observers" (when metrics are collected) and "collectors" (where/how the metrics are sent) are intended to be pluggable. The default, built-in collector is for Prometheus and will start a metrics server on the port passed to `natsby.NewPrometheusCollector`:

```go
collector := natsby.NewPrometheusCollector("2112")
observer := natsby.NewDefaultObserver(collector)
engine.Use(natsby.WithMetrics(observer))
```

## FAQ
**Why would I use this?**
If you find yourself needing to write a service and you usually reach for something like [Gin](https://github.com/gin-gonic/gin) but want to use [NATS](https://github.com/nats-io) instead of HTTP, Natsby might be for you.

**Why would I use NATS instead of HTTP?**
The [creators of NATS have more to say on this](https://changelog.com/gotime/130), but it really boils down to a few key things:

1. **Decoupling message passing and identity** - Services publishing messages don't need to know anything about who is handler those messages. No service discovery or DNS records to pass around. Services receiving those messages can change over time without any changes required on the sending side.
2. **Multiple messaging semantics** - Messages can be sent through your system in any combination of:
- consumed (intended for one receiver) or observed (intended for many receivers)
- Synchronous (sender expects a reply immediately) or asynchronous (sending accepts a reply later or never)

> NOTE: This terminology comes from [The Tao of Microservices](https://www.manning.com/books/the-tao-of-microservices) by Richard Roger.

You can replicate HTTP by using synchronous, consumed messages (i.e. request/reply), but any other combination is available. This can also change over time - two services might interact in a synchronous request/reply pattern, but you can also add other services that act as asynchronous observers. This is a totally additive change - nothing about the other collaborating services needs to change.

> NOTE: Out of the box, NATS does not support durable queues and only guarantees "at most once" delivery. This means that if there are no subscribers for a given topic, messages sent to it will never be received. Other projects such as [NATS Streaming Server](https://github.com/nats-io/nats-streaming-server), [Liftbridge](https://github.com/liftbridge-io/liftbridge) and [Jetstream](https://github.com/nats-io/jetstream), are building persistent messageing on top of NATS, but are out of scope of this project for now.

## Guiding Principles
- The API should feel familiar to developers familiar with [Gin](https://github.com/gin-gonic/gin), [Martini](https://github.com/go-martini/martini) or even [Express](https://expressjs.com/), [Sinatra](http://sinatrarb.com/), etc.
- Logic not core to NATS subscriptions/replies should live in middleware and be as pluggable as possible.
- Natsby should be able to do everything the [`nats-io` Go client](https://github.com/nats-io/nats.go) can do (this is a bit of a WIP).
- Observability is first-class.
