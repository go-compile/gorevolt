package gorevolt

import "errors"

func (c *Client) handleMessage(msg *message) {

	m := &Message{
		ID:        msg.ID,
		AuthorID:  msg.AuthorID,
		ChannelID: msg.ChannelID,
		Content:   msg.Content,

		c: c,
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

func (m *Message) Channel() (*Channel, error) {
	channel := m.c.cache.GetChannel(m.ChannelID)
	if channel != nil {
		return channel, nil
	}

	return nil, errors.New("could not fetch channel")
}
