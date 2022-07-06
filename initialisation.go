package gorevolt

import jsoniter "github.com/json-iterator/go"

func initialiseUserCache(c *Client) error {

	servers := c.cache.ListServers()
	for i := 0; i < len(servers); i++ {

		bucketServers.Wait()
		bucketServers.Increment(1)

		resp, err := c.request("GET", newRoute(RouteServerMembers, servers[i].ID), nil)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		// TODO: check response headers and if bucket out of sync, fix it

		var r serverMembers
		if err := jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
			return err
		}

		for x := 0; x < len(r.Users); x++ {
			c.cache.PutUser(r.Users[x])
		}
	}

	return nil
}
