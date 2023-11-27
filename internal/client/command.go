package client

import "github.com/willoma/keepakonf/internal/commands"

func (c *client) commands(a ...any) {
	callback(a, commands.List())
}
