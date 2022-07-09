package gorevolt

import "fmt"

// Server will provide the server object which the channel belongs to.
// If the server could not be found the result will be nil
func (channel *Channel) Server(c *Client) *Server {
	return c.cache.GetServer(channel.ServerID)
}

// SendMessage will send message in the channel
func (channel *Channel) SendMessage(c *Client, content string) (*Message, error) {
	return sendMessage(c, channel.ID, &newMessage{
		Content: content,
	})
}

// SendMessagef will format the message being sent to the channel
func (channel *Channel) SendMessagef(c *Client, format string, a ...interface{}) (*Message, error) {
	return sendMessage(c, channel.ID, &newMessage{
		Content: fmt.Sprintf(format, a...),
	})
}

func updateChannel(c *Client, update *channelUpdate) {

	old := c.cache.GetChannel(update.ID)
	if old == nil {
		// TODO: fetch channel and populate cache
		return
	}

	// clone channel
	current := *old

	for i := 0; i < len(update.Clear); i++ {
		switch update.Clear[i] {
		case "Icon":
			// TODO: channel updated icons
		case "Description":
			current.Description = ""
		}
	}

	// update changed fields
	for k, v := range update.Data {
		switch k {
		case "name":
			current.Name = v.(string)
		case "description":
			current.Description = v.(string)
		case "nfsw":
			current.NFSW = v.(bool)
		case "channel_type":
			current.ChannelType = v.(string)
		case "default_permissions":
			current.DefaultPermissions = interfacesToPermissions(v)
		case "role_permissions":
			current.RolePermissions = v.(map[string]Permissions)
		}
	}

	c.cache.PutChannel(&current)

	// Execute on channel updated handler
	for _, handler := range c.handlers.channelUpdate {
		go handler(c, old, &current)
	}
}

// Channels returns the channels which belong to this server
func (s *Server) Channels(c *Client) (channels []*Channel) {
	for i := 0; i < len(s.ChannelIDs); i++ {
		channels = append(channels, c.cache.GetChannel(s.ChannelIDs[i]))
	}

	return channels
}

func updateServer(c *Client, update *serverUpdate) {

	old := c.cache.GetServer(update.ID)
	if old == nil {
		// TODO: fetch server and populate cache
		return
	}

	// clone server
	current := *old

	for i := 0; i < len(update.Clear); i++ {
		switch update.Clear[i] {
		case "Icon":
			// TODO: server updated icons
		case "Description":
			current.Description = ""
		case "Banner":
			// TODO: clear server banner
		}
	}

	// update changed fields
	for k, v := range update.Data {
		switch k {
		case "name":
			current.Name = v.(string)
		case "description":
			current.Description = v.(string)
		case "owner":
			current.OwnerID = v.(string)
		case "channels":
			current.ChannelIDs = interfacesToStrings(v.([]interface{}))
		case "default_permissions":
			current.DefaultPermissions = v.(int64)
		case "roles":
			current.Roles = v.(map[string]Role)
		case "categories":
			current.Categories = interfacesToCategory(v.([]interface{}))
		}
	}

	c.cache.PutServer(&current)

	// Execute on channel updated handler
	for _, handler := range c.handlers.serverUpdate {
		go handler(c, old, &current)
	}
}

func interfacesToStrings(a []interface{}) []string {
	output := make([]string, len(a))

	for i := 0; i < len(a); i++ {
		output[i] = a[i].(string)
	}

	return output
}

func interfacesToCategory(a []interface{}) []Category {
	output := make([]Category, len(a))

	// preallocate a category map
	x := make(map[string]interface{}, 3)
	for i := 0; i < len(a); i++ {
		// cast type to map
		x = a[i].(map[string]interface{})

		// cast and place values into category
		output[i] = Category{
			ID:       x["id"].(string),
			Channels: interfacesToStrings(x["channels"].([]interface{})),
			Title:    x["title"].(string),
		}
	}

	return output
}

func interfacesToPermissions(a interface{}) Permissions {
	p := a.(map[string]interface{})

	return Permissions{
		A: int(p["a"].(float64)),
		D: int(p["d"].(float64)),
	}
}
