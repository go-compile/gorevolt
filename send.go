package gorevolt

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	jsoniter "github.com/json-iterator/go"
)

func (c *Client) SendMessage(channel, content string) (*Message, error) {

	_, err := sendMessage(c, channel, &newMessage{
		Content: content,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println(r.StatusCode)
	fmt.Println(string(body))

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return nil, errors.New("could not send message")
	}

	return nil, nil
}
