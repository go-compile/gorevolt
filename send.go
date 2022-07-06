package gorevolt

import (
	"bytes"
	"errors"

	jsoniter "github.com/json-iterator/go"
)

func (c *Client) SendMessage(channel, content string) (*Message, error) {

	m, err := sendMessage(c, channel, &newMessage{
		Content: content,
	})

	if err != nil {
		return nil, err
	}

	return m, nil
}

func sendMessage(c *Client, channel string, msg *newMessage) (*Message, error) {

	buf := bytes.NewBuffer(nil)
	if err := jsoniter.NewEncoder(buf).Encode(msg); err != nil {
		return nil, err
	}

	r, err := c.request("POST", newRoute(RouteChannelMessages, channel), buf)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return nil, errors.New("could not send message")
	}

	var m message
	if err := jsoniter.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}

	return c.convertMessage(&m), nil
}
