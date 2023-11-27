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

	instruction := c.data.GetInstruction(instructionID)
	if instruction == nil {
		return
	}
	instruction.Apply()
}
