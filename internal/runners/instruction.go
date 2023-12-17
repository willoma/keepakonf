package runners

import (
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

type Instruction interface {
	extractSaveable() map[string]any
	getInstruction(id string) (Instruction, bool)
	updateVariables(variables.Variables)
	getOutVariables() variables.Variables
	getStatus() status.Status
	watch()
	stop()
	Apply() bool
}

func instructionFromMap(iface any, vars map[string]string, grp *Group) Instruction {
	mapped, ok := iface.(map[string]any)
	if !ok {
		// TODO Send error response to caller
		return nil
	}

	insType, ok := mapped["type"].(string)
	if !ok {
		insType = "command"
	}

	switch insType {
	case "command":
		return instructionCommandFromMap(mapped, vars, grp)
	default:
		return nil
	}
}
