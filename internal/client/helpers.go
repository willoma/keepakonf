package client

import (
	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/log"
)

func (c *client) users(a ...any) {
	users, err := external.ListUsers()
	if err != nil {
		log.Error(err, "Could not list users")
	}
	callback(a, users)
}

func (c *client) globalVariables(a ...any) {
	callback(a, external.GlobalVariables)
}
