package client

func (c *client) applyGroup(a ...any) {
	if len(a) == 0 {
		return
	}

	groupID, ok := a[0].(string)
	if !ok {
		return
	}

	group := c.data.GetGroup(groupID)
	if group == nil {
		return
	}

	group.Apply()

}

func (c *client) applyInstruction(a ...any) {
	if len(a) == 0 {
		return
	}

	instructionID, ok := a[0].(string)
	if !ok {
		return
	}

	if instruction, ok := c.data.GetInstruction(instructionID); ok {
		instruction.Apply()
	}
}
