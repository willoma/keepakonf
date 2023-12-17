package runners

import (
	"bytes"
	"encoding/json"

	"github.com/rs/xid"
	"github.com/willoma/keepakonf/internal/commands"
	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

type instructionCommand struct {
	ID         string         `json:"id"`
	Command    string         `json:"command"`
	Parameters map[string]any `json:"parameters,omitempty"`

	Status status.Status   `json:"status"`
	Info   string          `json:"info"`
	Detail json.RawMessage `json:"detail"`

	command      commands.Command
	outVariables map[string]string

	group *Group
}

func (i *instructionCommand) extractSaveable() map[string]any {
	paramsClone := make(map[string]any, len(i.Parameters))
	for k, v := range i.Parameters {
		paramsClone[k] = v
	}
	return map[string]any{
		"id":         i.ID,
		"command":    i.Command,
		"parameters": paramsClone,
	}
}

func (i *instructionCommand) getInstruction(id string) (Instruction, bool) {
	if i.ID == id {
		return i, true
	}
	return nil, false
}

func (i *instructionCommand) updateVariables(vars variables.Variables) {
	i.command.UpdateVariables(vars)
}

func (i *instructionCommand) getOutVariables() variables.Variables {
	return i.outVariables
}

func (i *instructionCommand) getStatus() status.Status {
	return i.Status
}

func (i *instructionCommand) watch() {
	i.Status = status.StatusUnknown
	i.Info = "Checking..."
	i.Detail = nil
	i.command.Watch()
}

func (i *instructionCommand) stop() {
	if i.command != nil {
		i.command.Stop()
	}
}

func (i *instructionCommand) Apply() bool {
	if i.Status == status.StatusApplied {
		return true
	}
	if i.command != nil {
		return i.command.Apply()
	}
	return false
}

func (i *instructionCommand) updateStatus(newStatus status.Status, info string, detail status.Detail, outVars variables.Variables) {
	var detailJSON json.RawMessage
	if detail != nil {
		detailJSON = status.DetailJSON(detail)
	}
	storeAndEmit := func() {
		i.Status = newStatus
		i.Info = info
		i.Detail = detailJSON
		i.outVariables = outVars

		msg := map[string]any{
			"instruction": i.ID,
			"status":      newStatus,
			"info":        info,
		}
		if detail != nil {
			msg["detail"] = detailJSON
		}
		i.group.io.Emit("status", msg)
		i.group.updateStatusAndVariables()
	}

	desc := commands.GetDescription(i.Command)

	log := func() {
		log.Info(
			i.Command+": "+info,
			desc.Icon,
			newStatus, i.group.ID, i.ID, i.group.Name, detailJSON,
		)
	}

	if i.Status == status.StatusUnknown {
		storeAndEmit()
		return
	}

	switch newStatus {
	case i.Status:
		// When status does not change, only re-emit if info, detail and/or out vars change
		if info != i.Info {
			// and only log if info change
			log()
			storeAndEmit()
		} else if !bytes.Equal(detailJSON, i.Detail) {
			storeAndEmit()
		} else {
			var outvarsChanged bool
			if len(i.outVariables) == len(outVars) {
				outvarsChanged = true
			} else {
				for k, v := range i.outVariables {
					newValue, ok := outVars[k]
					if !ok || newValue != v {
						outvarsChanged = true
						break
					}
				}
			}
			if outvarsChanged {
				storeAndEmit()
			}
		}
	case status.StatusFailed:
		log()
		storeAndEmit()
	case status.StatusTodo:
		// Do not change status if it was failed and becomes todo, because it
		// means there was a recalculation by the watcher, but we want to keep
		// the failed status.
		if i.Status != status.StatusFailed {
			log()
			storeAndEmit()
		}
	case status.StatusRunning:
		log()
		storeAndEmit()
	case status.StatusApplied:
		log()
		storeAndEmit()
	}
}

func instructionCommandFromMap(mapped map[string]any, vars map[string]string, grp *Group) Instruction {
	id, ok := mapped["id"].(string)
	if !ok {
		id = xid.New().String()
	}

	command, _ := mapped["command"].(string)
	parameters, _ := mapped["parameters"].(map[string]any)

	i := &instructionCommand{
		ID:           id,
		Command:      command,
		Parameters:   parameters,
		outVariables: map[string]string{},
		group:        grp,
	}
	i.command = commands.Init(command, parameters, vars, i.updateStatus)
	return i
}
