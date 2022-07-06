package gorevolt

type responseHeader struct {
	Type string `json:"type"`
}

type authenticate struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type Ready struct {
	Users   []*User   `json:"users"`
	Servers []*Server `json:"servers"`
	// Channels []User   `json:"channels"`
}

type Server struct {
	ID          string `json:"_id"`
	OwnerID     string `json:"owner"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

type Message struct {
	ID        string `json:"_id"`
	AuthorID  string `json:"author"`
	ChannelID string `json:"channel"`
	Content   string `json:"content"`
}

func newAuthenticate(token string) authenticate {
	return authenticate{
		Type:  "Authenticate",
		Token: token,
	}
}
