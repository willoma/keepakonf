package client

import (
	"github.com/willoma/keepakonf/internal/runners"
)

func (c *client) groups(a ...any) {
	callback(a, c.data.GetGroups())
}

func (c *client) addGroup(a ...any) {
	if len(a) == 0 {
		return
	}

	grp := runners.GroupFromMap(a[0], c.io, c.logger)

	c.data.AppendGroup(grp)

	c.Broadcast().Emit("add group", grp)
	callback(a, grp)
}

func (c *client) modifyGroup(a ...any) {
	if len(a) == 0 {
		return
	}

	grp := runners.GroupFromMap(a[0], c.io, c.logger)

	if !c.data.ModifyGroup(grp) {
		return
	}

	c.Broadcast().Emit("modify group", grp)
	callback(a, grp)
}

func (c *client) removeGroup(a ...any) {
	if len(a) == 0 {
		return
	}

	grpID, ok := a[0].(string)
	if !ok {
		return
	}

	if !c.data.RemoveGroup(grpID) {
		return
	}

	c.Broadcast().Emit("remove group", grpID)
	callback(a, grpID)
}
