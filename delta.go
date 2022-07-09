package gorevolt

import (
	"errors"
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
