package natsby

import (
	"fmt"
	"log"
	"runtime/debug"
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *Context, err interface{})

// WithRecovery catches panics to prevent program crashes
func WithRecovery() HandlerFunc {
	return recovery(defaultRecoveryFunc)
}

// WithCustomRecovery catches panics with custom handler
func WithCustomRecovery(handle RecoveryFunc) HandlerFunc {
	return recovery(handle)
}

func defaultRecoveryFunc(c *Context, err interface{}) {
	c.Err = fmt.Errorf("Panic recovered: %v", err)
}

func recovery(handle RecoveryFunc) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.SetOutput(c.errWriter)
				log.Printf("panic recovered %v %s", err, debug.Stack())
				handle(c, err)
			}
		}()
		c.Next()
	}
}
