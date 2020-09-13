package natsby

// WithByteReply Checks for reply channel and sends back byte response
func WithByteReply() HandlerFunc {
	return func(c *Context) {
		c.Next()

		if c.Msg.Reply != "" {
			c.didReply = true
			c.Respond(c.ByteReplyPayload)
		}
	}
}

// WithJSONReply Checks for reply channel and sends back JSON response
func WithJSONReply() HandlerFunc {
	return func(c *Context) {
		c.Next()

		if c.EncodedConn != nil && c.Msg.Reply != "" {
			c.didReply = true
			c.EncodedConn.Publish(c.Msg.Reply, c.JSONReplyPayload)
		}
	}
}
