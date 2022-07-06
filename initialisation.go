package gorevolt

import jsoniter "github.com/json-iterator/go"

func initialiseUserCache(c *Client) error {

	servers := c.cache.ListServers()
	for i := 0; i < len(servers); i++ {

		// wait if we have met the rate limit
		bucketServers.Wait()
		bucketServers.Increment(1)

		resp, err := c.request("GET", newRoute(RouteServerMembers, servers[i].ID), nil)
		if err != nil {
			return err
		}

		// TODO: check response headers and if bucket out of sync, fix it

		var r serverMembers
		if err := jsoniter.NewDecoder(resp.Body).Decode(&r); err != nil {
			resp.Body.Close()
			return err
		}
		// close straight after decode instead of defer as defer would wait until
		// all servers are enumerated
		resp.Body.Close()

		for x := 0; x < len(r.Users); x++ {
			c.cache.PutUser(r.Users[x])
		}
	}

	return nil
}
