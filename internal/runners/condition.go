package runners

import "github.com/willoma/keepakonf/internal/variables"

type condition interface {
	extractSaveable() map[string]any
	check(vars variables.Variables) bool
}

type OpEqual struct {
	Left  string `json:"left"`
	Right string `json:"right"`
}

func (o OpEqual) extractSaveable() map[string]any {
	return map[string]any{
		"op":    "=",
		"left":  o.Left,
		"right": o.Right,
	}
}

func (o OpEqual) check(vars variables.Variables) bool {
	return vars.Replace(o.Left) == vars.Replace(o.Right)
}

func conditionFromMap(mapped map[string]any) (condition, bool) {
	op, ok := mapped["op"].(string)
	if !ok {
		return nil, false
	}

	switch op {
	case "=":
		return opEqualFromMap(mapped), true
	default:
		return nil, false
	}
}

func opEqualFromMap(mapped map[string]any) condition {
	left, _ := mapped["left"].(string)
	right, _ := mapped["right"].(string)
	return &OpEqual{left, right}
}
