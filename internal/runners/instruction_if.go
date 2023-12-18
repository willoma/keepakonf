package runners

import (
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

type instructionIf struct {
	Type         string        `json:"type"`
	Condition    condition     `json:"cond"`
	Instructions []Instruction `json:"instructions"`

	watching bool
	inVars   variables.Variables
}

func (i *instructionIf) extractSaveable() map[string]any {
	instructionsClone := make([]any, len(i.Instructions))
	for j, ins := range i.Instructions {
		instructionsClone[j] = ins.extractSaveable()
	}
	return map[string]any{
		"type":         "if",
		"cond":         i.Condition.extractSaveable(),
		"instructions": instructionsClone,
	}
}

func (i *instructionIf) getInstruction(id string) (Instruction, bool) {
	for _, ins := range i.Instructions {
		if iins, ok := ins.getInstruction(id); ok {
			return iins, true
		}
	}
	return nil, false
}

func (i *instructionIf) updateVariables(vars variables.Variables) {
	i.inVars = vars.Clone()
	localVars := vars.Clone()
	for _, ins := range i.Instructions {
		ins.updateVariables(localVars)
		for k, v := range ins.getOutVariables() {
			localVars[k] = v
		}
	}
}

func (i *instructionIf) getOutVariables() variables.Variables {
	vars := variables.Variables{}
	for _, ins := range i.Instructions {
		vars.Update(ins.getOutVariables())
	}
	return vars
}

func (i *instructionIf) getStatus() status.Status {
	var (
		instructionRunning bool
		instructionTodo    bool
		instructionFailed  bool
		instructionUnknown bool
	)

	for _, ins := range i.Instructions {
		switch ins.getStatus() {
		case status.StatusRunning:
			instructionRunning = true
		case status.StatusTodo:
			// Do not record as todo if there is a failed or unknown before
			if !instructionFailed && !instructionUnknown {
				instructionTodo = true
			}
		case status.StatusFailed:
			// Do not record as failed if there is a todo before
			if !instructionTodo {
				instructionFailed = true
			}
		case status.StatusUnknown:
			instructionUnknown = true
		}
	}

	switch {
	case instructionRunning:
		return status.StatusRunning
	case instructionFailed:
		return status.StatusFailed
	case instructionUnknown:
		return status.StatusUnknown
	case instructionTodo:
		return status.StatusTodo
	default:
		return status.StatusApplied
	}
}

func (i *instructionIf) watch() {
	if !i.check() {
		return
	}

	i.watching = true
	for _, ins := range i.Instructions {
		ins.watch()
	}
}

func (i *instructionIf) stop() {
	if i.watching {
		for _, ins := range i.Instructions {
			ins.stop()
		}
	}
	i.watching = false
}

func (i *instructionIf) Apply() bool {
	if !i.check() {
		return true
	}

	for _, ins := range i.Instructions {
		if !ins.Apply() {
			return false
		}
	}
	return true
}

func (i *instructionIf) check() bool {
	return i.Condition.check(i.inVars)
}

func instructionIfFromMap(mapped map[string]any, vars map[string]string, grp *Group) Instruction {
	condMap, ok := mapped["cond"].(map[string]any)
	if !ok {
		return nil
	}

	cond, ok := conditionFromMap(condMap)
	if !ok {
		return nil
	}

	instructionsIfaces, _ := mapped["instructions"].([]any)
	instructions := make([]Instruction, 0, len(instructionsIfaces))
	for _, ins := range instructionsIfaces {
		instruction := instructionFromMap(ins, vars, grp)
		if instruction != nil {
			instructions = append(instructions, instruction)
		}
	}

	i := &instructionIf{
		Type:         "if",
		Condition:    cond,
		Instructions: instructions,
	}
	return i
}
