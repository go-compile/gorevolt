package gorevolt

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	RevoltAPI       = "https://api.revolt.chat"
	RevoltWebsocket = "wss://ws.revolt.chat"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
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
	ws *websocket.Conn

	m sync.RWMutex
}

type handlers struct {
	ready []HandlerReady
}

type HandlerReady func(startup time.Duration)

// New creates a new client but does not authenticate yet
func New(token string) *Client {
	c := &Client{
		token:     token,
		api:       RevoltAPI,
		websocket: RevoltWebsocket,
		created:   time.Now(),
		http:      http.DefaultClient,
	}

	return c
}

// Register will setup a listener for your handler.
//
// Multiple handlers can be registered for the same event and
// be called concurrently.
func (c *Client) Register(handler interface{}) {
	c.m.Lock()
	defer c.m.Unlock()

	switch h := handler.(type) {
	case func(startup time.Duration):
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

// Close will disconnect you from Revolt
func (c *Client) Close() error {
	return c.ws.Close()
}

// Connect will establish a connection to the websocket server and
// authenticate your credentials
func (c *Client) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.websocket, nil)
	if err != nil {
		return err
	}

	if err := c.authenticate(conn); err != nil {
		return err
	}

	c.ws = conn

	if err := c.prepare(conn); err != nil {
		return err
	}

	// go handle events
	go c.eventLoop(conn)

	return nil
}

func (c *Client) authenticate(conn *websocket.Conn) error {
	// authenticate with websocket
	buf, err := jsoniter.Marshal(newAuthenticate(c.token))
	if err != nil {
		return err
	}

	if err := conn.WriteMessage(1, buf); err != nil {
		return err
	}

	_, buf, err = conn.ReadMessage()
	if err != nil {
		return err
	}

	var r responseHeader
	err = jsoniter.Unmarshal(buf, &r)
	if err != nil {
		return err
	}

	if r.Type != "Authenticated" {
		return ErrAuthenticationFailed
	}

	return nil
}

func (c *Client) prepare(conn *websocket.Conn) error {
	_, buf, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	var event Ready
	err = jsoniter.Unmarshal(buf, &event)
	if err != nil {
		return err
	}

	fmt.Println(event)

	// Execute ready handler
	for _, handler := range c.handlers.ready {
		go handler(time.Since(c.created))
	}

	return nil
}

func (c *Client) eventLoop(conn *websocket.Conn) {
	for {
		_, buf, err := conn.ReadMessage()
		if err != nil {
			// TODO: write to internal logger
			return
		}

		// decode header and find out event type
		var eventHeader responseHeader
		err = jsoniter.Unmarshal(buf, &eventHeader)
		if err != nil {
			// TODO: write to internal logger
			return
		}

		c.parseEvents(buf, eventHeader)
	}
}

func (c *Client) parseEvents(buf []byte, header responseHeader) {
	switch header.Type {

	case "Error":
		log.Printf("[WS:ERROR] %s\n", string(buf))
	case "Pong":
		// for now ignore
	case "Message":
		var msg Message
		err := jsoniter.Unmarshal(buf, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(msg)
	}

	fmt.Println(string(buf))
}
