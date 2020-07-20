package natsby

// WithByteReply Checks for reply channel and sends back byte response
func WithByteReply() HandlerFunc {
	return func(c *Context) {
		c.Next()

		if c.Msg.Reply != "" {
			c.didReply = true
			c.Engine.NatsConnection.Publish(c.Msg.Reply, c.ByteReplyPayload)
		}
	}
}

// WithJSONReply Checks for reply channel and sends back JSON response
func WithJSONReply() HandlerFunc {
	return func(c *Context) {
		c.Next()

		if c.Engine.NatsEncodedConnection != nil && c.Msg.Reply != "" {
			c.didReply = true
			c.Engine.NatsEncodedConnection.Publish(c.Msg.Reply, c.JSONReplyPayload)
		}
	}
}
