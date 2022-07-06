package gorevolt

import (
	"errors"
)

func (c *Client) handleMessage(msg *message) {

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

	for i := range c.handlers.message {
		c.handlers.message[i](c, m)
	}
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
