package gorevolt

import (
	"net/http"
	"sync"
	"time"
)

const (
	RevoltAPI       = "https://ws.revolt.chat"
	RevoltWebsocket = "wss://api.revolt.chat"
)

// Client maintains state for your logged in session
type Client struct {
	// http client used to send requests to the "Delta" API
	http *http.Client

	token string
	// is the server URL with the protocol
	// default: https://api.revolt.chat
	api       string
	websocket string
	created   time.Time

	handlers

	m sync.RWMutex
}

type handlers struct {
	ready []HandlerReady
}

type HandlerReady func(startup time.Duration)

// New creates a new client but does not authenticate yet
func New(token string) *Client {
	c := &Client{
		token:   token,
		api:     RevoltAPI,
		created: time.Now(),
		http:    http.DefaultClient,
	}

	return c
}

// RegisterHandler will setup a listener for your specified event.
//
// Multiple handlers can be registered for the same event and
// be called concurrently.
func (c *Client) RegisterHandler(callback interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	switch h := callback.(type) {
	case HandlerReady:
		c.handlers.ready = append(c.handlers.ready, h)
	default:
		panic("unknown handler")
	}
}

// UnregisterHandlers will reset all handlers and require
// new ones to be set.
func (c *Client) UnregisterHandlers() {
	c.m.Lock()
	defer c.m.Unlock()

	c.handlers = handlers{}
}

func (c *Client) Connect() error {
	return nil
}
