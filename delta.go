package gorevolt

import (
	"errors"

	jsoniter "github.com/json-iterator/go"
)

var (
	ErrRateLimit    = errors.New("429: rate limit reached")
	ErrUnauthorised = errors.New("unauthorised request")
	ErrForbidden    = errors.New("access denied you do not have permission to perform that action")
	ErrUnknown      = errors.New("unknown error occurred")
	ErrNotFound     = errors.New("resource not found")
)

// parseStatus takes in a non success status code and returns the
// appropriate error
func parseStatus(status int) error {

	switch status {
	case 429:
		return ErrRateLimit
	case 401:
		return ErrUnauthorised
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	}

	return ErrUnknown
}

func GetChannel(c *Client, id string) (*Channel, error) {
	r, err := c.request("GET", newRoute(RouteChannel, id), nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return nil, parseStatus(r.StatusCode)
	}

	var m Channel
	if err := jsoniter.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

func GetUser(c *Client, id string) (*User, error) {
	r, err := c.request("GET", newRoute(RouteUsers, id), nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		return nil, parseStatus(r.StatusCode)
	}

	var m User
	if err := jsoniter.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}

	return &m, nil
}
