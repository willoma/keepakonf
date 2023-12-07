package runners

import (
	"encoding/json"

	"github.com/rs/xid"

	"github.com/willoma/keepakonf/internal/commands"
	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
)

type Instruction struct {
	ID         string         `json:"id"`
	Command    string         `json:"command"`
	Parameters map[string]any `json:"parameters,omitempty"`

	Status status.Status   `json:"status"`
	Info   string          `json:"info"`
	Detail json.RawMessage `json:"detail"`

	command commands.Command

	group *Group
}

func (i *Instruction) watch() {
	i.Status = status.StatusUnknown
	i.Info = "Checking..."
	i.Detail = nil
	i.command.Watch()
}

func (i *Instruction) stop() {
	if i.command != nil {
		i.command.Stop()
	}
}

func (i *Instruction) Apply() bool {
	if i.command != nil {
		return i.command.Apply()
	}
	return false
}

func (i *Instruction) updateStatus(newStatus status.Status, info string, detail status.Detail) {
	var detailJSON json.RawMessage
	if detail != nil {
		detailJSON = detail.JSON()
	}
	storeAndEmit := func() {
		i.Status = newStatus
		i.Info = info
		i.Detail = detailJSON

		msg := map[string]any{
			"instruction": i.ID,
			"status":      newStatus,
			"info":        info,
		}
		if detail != nil {
			msg["detail"] = detailJSON
		}
		i.group.io.Emit("status", msg)
		i.group.updateStatus()
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
		// When status does not change, log only if info changes
		if info != i.Info {
			log()
		}
		storeAndEmit()
	case status.StatusFailed:
		log()
		storeAndEmit()
	case status.StatusNone:
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

func instructionFromMap(iface any, grp *Group) *Instruction {
	mapped, ok := iface.(map[string]any)
	if !ok {
		// TODO Send error response to caller
		return nil
	}

	id, ok := mapped["id"].(string)
	if !ok {
		id = xid.New().String()
	}

	command, _ := mapped["command"].(string)
	parameters, _ := mapped["parameters"].(map[string]any)

	i := &Instruction{
		ID:         id,
		Command:    command,
		Parameters: parameters,
		group:      grp,
	}
	i.command = commands.Init(command, parameters, i.updateStatus)
	return i
}
