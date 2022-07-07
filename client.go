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

	pingInterval = 15 * time.Second
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
	// wsM mutex for writing to the websocket
	wsM     sync.Mutex
	wsClose <-chan struct{}

	User  *User
	cache Cache

	m sync.RWMutex
}

type handlers struct {
	ready         []HandlerReady
	message       []HandlerMessage
	messageUpdate []HandlerMessageUpdate
}

type HandlerReady func(c *Client, startup time.Duration)
type HandlerMessage func(c *Client, m *Message)
type HandlerMessageUpdate func(c *Client, m *UpdatedMessage)

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

// OnReady registers a ready event handler
func (c *Client) OnReady(h HandlerReady) {
	c.handlers.ready = append(c.handlers.ready, h)
}

// OnMessage registers a onMessage event handler
func (c *Client) OnMessage(h HandlerMessage) {
	c.handlers.message = append(c.handlers.message, h)
}

// OnMessage registers a onMessageUpdate event handler
func (c *Client) OnMessageUpdate(h HandlerMessageUpdate) {
	c.handlers.messageUpdate = append(c.handlers.messageUpdate, h)
}

// SetCache allows you to use custom caching layers.
// Solutions such as hash maps, Redis or even disk
// caches are possible.
func (c *Client) SetCache(cache Cache) {
	c.cache = cache
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
	if c.cache == nil {
		c.cache = NewArrayCache(300, 20)
	}

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
	go c.pingLoop(conn)

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

	for i := range event.Users {
		// First user is the current user
		if i == 0 {
			c.User = event.Users[i]
		}

		c.cache.PutUser(event.Users[i])
	}

	for i := range event.Servers {
		c.cache.PutServer(event.Servers[i])
	}

	for i := range event.Channels {
		c.cache.PutChannel(event.Channels[i])
	}

	if err = initialiseUserCache(c); err != nil {
		return err
	}

	// Execute ready handler
	for _, handler := range c.handlers.ready {
		go handler(c, time.Since(c.created))
	}

	return nil
}

func (c *Client) eventLoop(conn *websocket.Conn) {
	for {
		// TODO: send ping events
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

func (c *Client) pingLoop(conn *websocket.Conn) {
	t := time.NewTicker(pingInterval)

	for {
		select {
		case <-t.C:
			c.wsM.Lock()
			// pingBuf is precomputed for optimum performance
			conn.WriteMessage(1, pingBuf)
			c.wsM.Unlock()
		case <-c.wsClose:
			t.Stop()
		}
	}
}

func (c *Client) parseEvents(buf []byte, header responseHeader) {
	switch header.Type {

	case "Error":
		log.Printf("[WS:ERROR] %s\n", string(buf))
	case "Pong":
		// for now ignore
	case "Message":
		var msg message
		err := jsoniter.Unmarshal(buf, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		c.handleMessage(&msg)
	case "MessageUpdate":
		var msg MessageUpdate
		err := jsoniter.Unmarshal(buf, &msg)
		if err != nil {
			log.Println(err)
			return
		}

		c.handleUpdatedMessage(&msg)
	}

	fmt.Println(string(buf))
}
