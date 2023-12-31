package runners

import (
	"github.com/rs/xid"
	"github.com/zishang520/socket.io/v2/socket"

	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`

	Instructions []Instruction `json:"instructions"`

	Status status.Status `json:"status"`

	io socket.NamespaceInterface
}

func (g *Group) ExtractSaveable() map[string]any {
	instructionsClone := make([]any, len(g.Instructions))
	for i, ins := range g.Instructions {
		instructionsClone[i] = ins.extractSaveable()
	}
	return map[string]any{
		"id":           g.ID,
		"name":         g.Name,
		"icon":         g.Icon,
		"instructions": instructionsClone,
	}
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
	for _, ins := range g.Instructions {
		if !ins.Apply() {
			return
		}
	}
}

func (g *Group) GetInstruction(id string) (Instruction, bool) {
	for _, i := range g.Instructions {
		if ins, ok := i.getInstruction(id); ok {
			return ins, true
		}
	}
	return nil, false
}

func (g *Group) updateStatusAndVariables() {
	vars := variables.GlobalMap()

	var (
		instructionRunning bool
		instructionTodo    bool
		instructionFailed  bool
		instructionUnknown bool
	)

	for _, ins := range g.Instructions {
		switch ins.getStatus() {
		case status.StatusRunning:
			instructionRunning = true
		case status.StatusTodo:
			// Do not record the group as todo if there is a failed or unknown before
			if !instructionFailed && !instructionUnknown {
				instructionTodo = true
			}
		case status.StatusFailed:
			// Do not record the group as failed if there is a todo before
			if !instructionTodo {
				instructionFailed = true
			}
		case status.StatusUnknown:
			instructionUnknown = true
		}

		ins.updateVariables(vars)
		for k, v := range ins.getOutVariables() {
			vars[k] = v
		}
	}

	var newStatus status.Status
	switch {
	case instructionRunning:
		newStatus = status.StatusRunning
	case instructionFailed:
		newStatus = status.StatusFailed
	case instructionUnknown:
		newStatus = status.StatusUnknown
	case instructionTodo:
		newStatus = status.StatusTodo
	default:
		newStatus = status.StatusApplied
	}

	if g.Status != newStatus {
		g.Status = newStatus
		g.io.Emit("group status", map[string]any{
			"group":  g.ID,
			"status": newStatus,
		})
	}
}

func GroupFromMap(iface any, io socket.NamespaceInterface) *Group {
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
	icon, _ := mapped["icon"].(string)
	grp := &Group{
		ID:     id,
		Name:   name,
		Icon:   icon,
		Status: status.StatusUnknown,
		io:     io,
	}

	vars := variables.GlobalMap()

	instructionsIfaces, _ := mapped["instructions"].([]any)
	instructions := make([]Instruction, 0, len(instructionsIfaces))
	for _, ins := range instructionsIfaces {
		instruction := instructionFromMap(ins, vars, grp)
		if instruction != nil {
			instructions = append(instructions, instruction)
		}
	}

	grp.Instructions = instructions

	return grp
}
