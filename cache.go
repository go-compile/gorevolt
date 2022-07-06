package gorevolt

import "sync"

type Cache interface {
	GetUser(id string) *User
	PutUser(user *User)
	GetServer(id string) *Server
	PutServer(server *Server)
	GetChannel(id string) *Channel
	PutChannel(channel *Channel)
}

const avgChannelsPerServer = 10

// ArrayCache uses a array to store users, servers etc.
//
// Results in a reduced memory footprint but slower
// performance.
type ArrayCache struct {
	users    []*User
	servers  []*Server
	channels []*Channel

	m sync.RWMutex
}

// NewArrayCache creates a new cache layer using minimal memory.
//
// user and server hint is used to preallocate memory to boost
// performance. It is better to be over generous than conservative.
func NewArrayCache(usersHint, serversHint int) Cache {
	return &ArrayCache{
		users:    make([]*User, 0, usersHint),
		servers:  make([]*Server, 0, serversHint),
		channels: make([]*Channel, 0, serversHint*avgChannelsPerServer),
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
		// overwrite existing server
		server = s
		return
	}

	c.servers = append(c.servers, s)
}

func (c *ArrayCache) GetChannel(id string) *Channel {
	c.m.RLock()
	defer c.m.RUnlock()

	return c.getChannel(id)
}

func (c *ArrayCache) getChannel(id string) *Channel {

	for i := 0; i < len(c.channels); i++ {
		if c.channels[i].ID == id {
			return c.channels[i]
		}
	}

	return nil
}

func (c *ArrayCache) PutChannel(channel *Channel) {
	c.m.Lock()
	defer c.m.Unlock()

	if chan1 := c.getChannel(channel.ID); chan1 != nil {
		// overwrite existing channel
		chan1 = channel
		return
	}

	c.channels = append(c.channels, channel)
}
