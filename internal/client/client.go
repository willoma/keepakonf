package client

import (
	"github.com/zishang520/socket.io/v2/socket"

	"github.com/willoma/keepakonf/internal/data"
)

type client struct {
	*socket.Socket
	io   socket.NamespaceInterface
	data *data.Data
}

func Serve(s *socket.Socket, io socket.NamespaceInterface, data *data.Data) {
	c := client{
		Socket: s,
		io:     io,
		data:   data,
	}

	c.On("apply group", c.applyGroup)
	c.On("apply instruction", c.applyInstruction)

	c.On("commands", c.commands)

	c.On("groups", c.groups)
	c.On("add group", c.addGroup)
	c.On("modify group", c.modifyGroup)
	c.On("remove group", c.removeGroup)

	c.On("logs", c.logs)

	c.On("users", c.users)
	c.On("global variables", c.globalVariables)
}

func callback(request []any, response ...any) {
	if len(request) == 0 {
		return
	}
	c, ok := request[len(request)-1].(func([]any, error))
	if !ok {
		return
	}
	c(response, nil)
}
