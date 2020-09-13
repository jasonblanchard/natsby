package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jasonblanchard/natsby"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// PingController controller for handling pings
type PingController struct {
	Logger zerolog.Logger
}

// Handle handles the event
func (ctrl *PingController) Handle(e *natsby.Event) (natsby.EventResult, error) {
	ctrl.Logger.Info().Msg(fmt.Sprintf("Recived message from subject: %s", e.Msg.Subject))
	if e.Msg.Reply != "" {
		e.Msg.Respond([]byte("pong"))
	}

	return natsby.EventResult{}, nil
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.DurationFieldUnit = time.Second

	e := &natsby.Engine{
		Conn:   nc,
		Logger: &logger,
	}

	pingController := &PingController{
		Logger: logger,
	}

	e.Subscribe("ping", pingController)

	err = e.Run(func() {
		e.Logger.Info().Msg("Ready 🚀")
	})

	if err != nil {
		panic(err)
	}
}
