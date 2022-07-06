package gorevolt

import "sync"

type Cache interface {
	GetUser(id string) *User
	PutUser(user *User)
	GetServer(id string) *Server
	PutServer(server *Server)
}

// ArrayCache uses a array to store users, servers etc.
//
// Results in a reduced memory footprint but slower
// performance.
type ArrayCache struct {
	users   []*User
	servers []*Server

	m sync.RWMutex
}

// NewArrayCache creates a new cache layer using minimal memory.
//
// user and server hint is used to preallocate memory to boost
// performance. It is better to be over generous than conservative.
func NewArrayCache(usersHint, serversHint int) Cache {
	return &ArrayCache{
		users:   make([]*User, 0, usersHint),
		servers: make([]*Server, 0, serversHint),
	}
}

func (c *ArrayCache) GetUser(id string) *User {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.getUser(id)
}
func (c *ArrayCache) getUser(id string) *User {
	for i := 0; i < len(c.users); i++ {
		if c.users[i].ID == id {
			return c.users[i]
		}
	}

	return nil
}

func (c *ArrayCache) GetServer(id string) *Server {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.getServer(id)
}

func (c *ArrayCache) getServer(id string) *Server {

	for i := 0; i < len(c.servers); i++ {
		if c.servers[i].ID == id {
			return c.servers[i]
		}
	}

	return nil
}

func (c *ArrayCache) PutUser(u *User) {
	c.m.Lock()
	defer c.m.Unlock()

	if user := c.getUser(u.ID); user != nil {
		// overwrite existing user
		user = u
		return
	}

	c.users = append(c.users, u)
}

func (c *ArrayCache) PutServer(s *Server) {
	c.m.Lock()
	defer c.m.Unlock()

	if server := c.getServer(s.ID); server != nil {
		// overwrite existing user
		server = s
		return
	}

	c.servers = append(c.servers, s)
}
