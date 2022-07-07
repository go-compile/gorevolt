package gorevolt

import "time"

var pingBuf = []byte(`{"type":"Ping","data":0}`)

type responseHeader struct {
	Type string `json:"type"`
}

type authenticate struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type Ready struct {
	Users    []*User    `json:"users"`
	Servers  []*Server  `json:"servers"`
	Channels []*Channel `json:"channels"`
}

type Server struct {
	ID                 string          `json:"_id"`
	OwnerID            string          `json:"owner"`
	Name               string          `json:"name"`
	Description        string          `json:"description"`
	Channels           []string        `json:"channels"`
	Categories         []Category      `json:"categories"`
	Roles              map[string]Role `json:"roles"`
	DefaultPermissions int64           `json:"default_permissions"`
}

type ServerCreate struct {
	ID       string    `json:"_id"`
	Server   Server    `json:"server"`
	Channels []Channel `json:"channels"`
}

type Channel struct {
	ID                 string                 `json:"_id"`
	ChannelType        string                 `json:"channel_type"`
	ServerID           string                 `json:"server"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	NFSW               bool                   `json:"nsfw"`
	DefaultPermissions Permissions            `json:"default_permissions"`
	RolePermissions    map[string]Permissions `json:"role_permissions"`
	LastMessageID      string                 `json:"last_message_id"`
}

type Category struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Channels []string `json:"channels"`
}

type Role struct {
	Name        string      `json:"name"`
	Permissions Permissions `json:"permissions"`
	Rank        int         `json:"rank"`
}

type Permissions struct {
	A int `json:"a"`
	D int `json:"d"`
}

type User struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
	Profile  struct {
		Content string `json:"content"`
	} `json:"profile"`
	Bot *struct {
		OwnerID string `json:"owner"`
	} `json:"bot"`
	Relationship string  `json:"relationship"`
	Online       bool    `json:"online"`
	Privileged   bool    `json:"privileged"`
	Flags        []int32 `json:"flags"`
}

type Member struct {
	ID struct {
		Server string `json:"server"`
		User   string `json:"user"`
	} `json:"_id"`
	Nickname string `json:"nickname"`
	// TODO: add avatar to member
	Roles []string `json:"roles"`
}

type UpdatedMessage struct {
	ID        string    `json:"id"`
	ChannelID string    `json:"channel_id"`
	ServerID  string    `json:"server_id"`
	Channel   *Channel  `json:"channel"`
	Content   string    `json:"content"`
	Edited    time.Time `json:"edited"`
	Changes   []string  `json:"changes"`
}

type MessageUpdate struct {
	ID      string                 `json:"id"`
	Channel string                 `json:"channel"`
	Data    map[string]interface{} `json:"data"`
}

type UpdatedChannel struct {
	ID    string   `json:"id"`
	Data  Channel  `json:"data"`
	Clear []string `json:"clear"`
}

type channelUpdate struct {
	ID    string                 `json:"id"`
	Data  map[string]interface{} `json:"data"`
	Clear []string               `json:"clear"`
}

type message struct {
	ID        string `json:"_id"`
	AuthorID  string `json:"author"`
	ChannelID string `json:"channel"`
	Content   string `json:"content"`
	Edited    string `json:"edited"`
}

type Message struct {
	ID        string   `json:"_id"`
	AuthorID  string   `json:"author_id"`
	ChannelID string   `json:"channel_id"`
	ServerID  string   `json:"server_id"`
	Channel   *Channel `json:"channel"`
	Author    *User    `json:"author"`
	Content   string   `json:"content"`

	c *Client
}

// newMessage is different to message as it has additional options to masquerade
type newMessage struct {
	Content     string      `json:"content"`
	Attachments []string    `json:"attachments"`
	Replies     []Reply     `json:"replies"`
	Embed       []Embed     `json:"embeds"`
	Masquerade  *Masquerade `json:"masquerade"`
}

type Reply struct {
	ID      string `json:"id"`
	Mention bool   `json:"mention"`
}

type Embed struct {
	IconURL     string `json:"icon_url"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Media       string `json:"media"`
	Colour      string `json:"colour"`
}

type Masquerade struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type serverMembers struct {
	Members []*Member
	Users   []*User
}

func newAuthenticate(token string) authenticate {
	return authenticate{
		Type:  "Authenticate",
		Token: token,
	}
}
