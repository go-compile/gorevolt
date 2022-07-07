package gorevolt

import (
	"fmt"
	"time"
)

func (c *Client) handleMessage(msg *message) {
	m := c.convertMessage(msg)

	for i := range c.handlers.message {
		c.handlers.message[i](c, m)
	}
}

func (c *Client) handleUpdatedMessage(msg *MessageUpdate) {
	m := &UpdatedMessage{
		ID:        msg.ID,
		ChannelID: msg.Channel,
		Changes:   []string{},
	}

	// TODO: add embed support
	content, ok := msg.Data["content"].(string)
	if ok {
		m.Content = content
		m.Changes = append(m.Changes, "content")
	}

	edited, ok := msg.Data["edited"].(string)
	if ok {
		m.Changes = append(m.Changes, "edited")
		stamp, err := time.Parse(time.RFC3339, edited)
		if err == nil {
			m.Edited = stamp
		}
	}

	for i := range c.handlers.messageUpdate {
		c.handlers.messageUpdate[i](c, m)
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

	author := m.c.cache.GetUser(m.AuthorID)
	if author != nil {
		m.Author = author
	}

	return m
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
