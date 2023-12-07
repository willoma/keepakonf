package commands

import (
	"github.com/willoma/keepakonf/internal/status"
)

type command struct {
	msg status.SendStatus
}

type Command interface {
	Watch()
	Stop()
	Apply() bool
}

type constructor func(params map[string]any, msg status.SendStatus) Command

type definition struct {
	description Description
	constructor constructor
}

var (
	byName = map[string]definition{}
)

func register(name, icon, description string, parameters ParamsDesc, c constructor) struct{} {
	byName[name] = definition{Description{name, icon, description, parameters}, c}
	return struct{}{}
}

func Init(name string, params map[string]any, msg status.SendStatus) Command {
	def, ok := byName[name]
	if !ok {
		return nil
	}
	def.description.Parameters.EnsureTyped(params)
	return def.constructor(params, msg)
}
