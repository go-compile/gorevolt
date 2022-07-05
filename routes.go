package gorevolt

type Route struct {
	route     string
	ratelimit int
}

const (
	RouteUsersMe = "/users/@me"
)
