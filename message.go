package gorevolt

import (
	"errors"
	"fmt"
)

func (c *Client) handleMessage(msg *message) {
	m := c.convertMessage(msg)

	for i := range c.handlers.message {
		c.handlers.message[i](c, m)
	}
}

// convertMessage takes a raw message received from the API or websocket
// and converts it to a user facing Message
func (c *Client) convertMessage(msg *message) *Message {
	m := &Message{
		ID:        msg.ID,
		AuthorID:  msg.AuthorID,
		ChannelID: msg.ChannelID,
		Content:   msg.Content,

		c: c,
	}

	channel := m.c.cache.GetChannel(m.ChannelID)
	if channel != nil {
		m.Channel = channel
		m.ServerID = channel.ServerID
	}

	return m
}

func (m *Message) Author() (*User, error) {
	user := m.c.cache.GetUser(m.AuthorID)
	if user != nil {
		return user, nil
	}

	return nil, errors.New("could not fetch author")
}

func (m *Message) Server() *Server {
	s := m.c.cache.GetServer(m.ServerID)
	if s != nil {
		return s
	}

	return nil
}

// Reply will send a message within the same channel with that message as the
// reply.
func (m *Message) Reply(content string) (*Message, error) {
	return sendMessage(m.c, m.ChannelID, &newMessage{
		Content: content,
		Replies: []Reply{{m.ID, true}},
	})
}

// Replyf is much like Message.Reply() however you can format your strings
// just like fmt.Printf()
func (m *Message) Replyf(format string, a ...interface{}) (*Message, error) {
	return sendMessage(m.c, m.ChannelID, &newMessage{
		Content: fmt.Sprintf(format, a...),
		Replies: []Reply{{m.ID, true}},
	})
}
