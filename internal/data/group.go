package data

import (
	"fmt"
	"slices"

	"github.com/willoma/keepakonf/internal/runners"
)

func (d *Data) GetGroups() []*runners.Group {
	d.mu.Lock()
	defer d.mu.Unlock()

	return d.groups
}

func (d *Data) GetGroup(id string) *runners.Group {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, g := range d.groups {
		if g.ID == id {
			return g
		}
	}

	return nil
}

func (d *Data) AppendGroup(group *runners.Group) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.groups = append(d.groups, group)
	d.save()
	d.logger.Info(
		fmt.Sprintf("Added group %q", group.Name),
		"group",
		"", group.ID, "", "", nil,
	)
	group.Watch()
}

func (d *Data) ModifyGroup(group *runners.Group) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	i := slices.IndexFunc(d.groups, func(grp *runners.Group) bool {
		return grp.ID == group.ID
	})
	if i == -1 {
		return false
	}
	d.groups[i].StopWatch()
	d.groups[i] = group
	d.save()
	d.logger.Info(
		fmt.Sprintf("Modified group %q", group.Name),
		"group",
		"", group.ID, "", "", nil,
	)
	group.Watch()
	return true
}

func (d *Data) RemoveGroup(id string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	var name string
	d.groups = slices.DeleteFunc(d.groups, func(grp *runners.Group) bool {
		if grp.ID == id {
			name = grp.Name
			grp.StopWatch()
			return true
		}
		return false
	})
	if name == "" {
		return true
	}
	d.save()
	d.logger.Info(
		fmt.Sprintf("Removed group %q", name),
		"group",
		"", "", "", "", nil,
	)
	return true
}
