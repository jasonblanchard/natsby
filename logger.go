package natsby

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// DefaultLogger sets up a simple zerolog instance
func DefaultLogger() *zerolog.Logger {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	zerolog.DurationFieldUnit = time.Second
	return &logger
}

// WithLogger wraps handler with logging
func WithLogger(logger *zerolog.Logger) HandlerFunc {
	return func(c *Context) {
		logger.Debug().
			Str("subject", c.Msg.Subject).
			Msg("received")

		latency := c.NextWithLatencyDuration()

		if c.Err != nil {
			logger.Error().
				Str("subject", c.Msg.Subject).
				Dur("latency", latency).
				Str("replyChan", c.Msg.Reply).
				Err(c.Err).
				Msg(fmt.Sprintf("%+v", c.Err))
			return
		}

		logger.Info().
			Str("subject", c.Msg.Subject).
			Dur("latency", latency).
			Str("replyChan", c.Msg.Reply).
			Msg("processed")
	}
}
