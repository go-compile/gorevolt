package gorevolt

import (
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/valyala/fasttemplate"
)

type Route struct {
	route     string
	ratelimit int
}

const (
	RouteUsersMe       = "/users/@me"
	RouteServerMembers = "/servers/{0}/members"
)

// newRoute takes in a existing route then inputs the params to the URL
func newRoute(route string, params ...string) string {
	t := fasttemplate.New(route, "{", "}")
	s := t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		i, err := strconv.Atoi(tag)
		if err != nil {
			panic("invalid tag type in newRoute")
		}

		return w.Write([]byte(url.PathEscape(params[i])))
	})

	return s
}

func (c *Client) request(method, path string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest(method, c.api+path, body)
	if err != nil {
		return nil, err
	}

	r.Header.Set("x-bot-token", c.token)

	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
