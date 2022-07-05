package gorevolt

type responseHeader struct {
	Type string `json:"type"`
}

type authenticate struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func newAuthenticate(token string) authenticate {
	return authenticate{
		Type:  "Authenticate",
		Token: token,
	}
}
