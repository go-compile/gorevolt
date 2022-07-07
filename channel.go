package gorevolt

import "fmt"

// Server will provide the server object which the channel belongs to.
// If the server could not be found the result will be nil
func (channel *Channel) Server(c *Client) *Server {
	return c.cache.GetServer(channel.ServerID)
}

// SendMessage will send message in the channel
func (channel *Channel) SendMessage(c *Client, content string) (*Message, error) {
	return sendMessage(c, channel.ID, &newMessage{
		Content: content,
	})
}

// SendMessagef will format the message being sent to the channel
func (channel *Channel) SendMessagef(c *Client, format string, a ...interface{}) (*Message, error) {
	return sendMessage(c, channel.ID, &newMessage{
		Content: fmt.Sprintf(format, a...),
	})
}
