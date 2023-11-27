package data

import "github.com/willoma/keepakonf/internal/runners"

func (d *Data) GetInstruction(id string) *runners.Instruction {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, grp := range d.groups {
		for _, ins := range grp.Instructions {
			if id == ins.ID {
				return ins
			}
		}
	}
	return nil
}
