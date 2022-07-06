package gorevolt

import (
	"io"
	"net/http"
)

type Route struct {
	route     string
	ratelimit int
}

const (
	RouteUsersMe = "/users/@me"
)

func (c *Client) request(method, path string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, c.api+path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
