package data

import "github.com/willoma/keepakonf/internal/runners"

func (d *Data) GetInstruction(id string) (runners.Instruction, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, grp := range d.groups {
		if ins, ok := grp.GetInstruction(id); ok {
			return ins, true
		}
	}
	return nil, false
}
