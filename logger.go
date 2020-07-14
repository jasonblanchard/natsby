package natsby

import (
	"fmt"
	"time"
)

// WithLogger wraps handler with logging
func WithLogger() HandlerFunc {
	return func(c *Context) {
		if c.Engine.Logger == nil {
			c.Next()
			return
		}

		c.Engine.Logger.Debug().
			Str("subject", c.Msg.Subject).
			Msg("received")

		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if c.Err != nil {
			c.Engine.Logger.Error().
				Str("subject", c.Msg.Subject).
				Err(c.Err).
				Msg(fmt.Sprintf("%+v", c.Err))
			return
		}

		c.Engine.Logger.Info().
			Str("subject", c.Msg.Subject).
			Dur("latencyMS", latency).
			Str("replyChan", c.Msg.Reply).
			Msg("processed")
	}
}
