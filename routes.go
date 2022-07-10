package gorevolt

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/valyala/fasttemplate"
)

const (
	RouteUsersMe         = "/users/@me"
	RouteUsers           = "/users/{0}"
	RouteServerMembers   = "/servers/{0}/members"
	RouteChannelMessages = "/channels/{0}/messages"
	RouteChannel         = "/channels/{0}"
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

	nonce, err := newNonce()
	if err != nil {
		return nil, err
	}

	r.Header.Set("x-bot-token", c.token)
	r.Header.Set("Idempotency-Key", nonce)

	resp, err := c.http.Do(r)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func newNonce() (string, error) {
	// TODO: make a faster unique ID/nonce generator
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", nil
	}

	return hex.EncodeToString(buf), nil
}
