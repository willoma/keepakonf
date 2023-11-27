package runners

import (
	"github.com/rs/xid"
	"github.com/zishang520/socket.io/v2/socket"

	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
)

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Instructions []*Instruction `json:"instructions"`

	Status status.Status `json:"status"`

	io     socket.NamespaceInterface
	logger *log.Logger
}

func (g *Group) Watch() {
	for _, i := range g.Instructions {
		i.watch()
	}
}

func (g *Group) StopWatch() {
	for _, i := range g.Instructions {
		i.stop()
	}
}

func (g *Group) Apply() {
	for _, i := range g.Instructions {
		if i.Status == status.StatusApplied {
			continue
		}
		if !i.Apply() {
			return
		}
	}

}

func (g *Group) updateStatus() {
	status := status.StatusApplied
	for _, ins := range g.Instructions {
		status = status.IfHigherPriority(ins.Status)
	}

	if g.Status != status {
		g.Status = status
		g.io.Emit("group status", map[string]any{
			"group":  g.ID,
			"status": status,
		})
	}
}

func GroupFromMap(iface any, io socket.NamespaceInterface, logger *log.Logger) *Group {
	mapped, ok := iface.(map[string]any)
	if !ok {
		// TODO Send error response to caller
		return nil
	}

	id, ok := mapped["id"].(string)
	if !ok {
		id = xid.New().String()
	}

	name, _ := mapped["name"].(string)

	grp := &Group{
		ID:     id,
		Name:   name,
		Status: status.StatusUnknown,
		io:     io,
		logger: logger,
	}

	instructionsIfaces, _ := mapped["instructions"].([]any)
	instructions := make([]*Instruction, len(instructionsIfaces))
	for i, ins := range instructionsIfaces {
		instruction := instructionFromMap(ins, grp)
		instructions[i] = instruction
	}

	grp.Instructions = instructions

	return grp
}
