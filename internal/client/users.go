package client

import "github.com/willoma/keepakonf/internal/external"

func (c *client) users(a ...any) {
	users, err := external.ListUsers()
	if err != nil {
		c.logger.Error(err, "Could not list users")
	}
	callback(a, users)
}
