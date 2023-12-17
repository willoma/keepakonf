package status

import "github.com/willoma/keepakonf/internal/variables"

type Status string

const (
	// StatusNone may be used in details (eg. cell color), but not in messages
	StatusNone    Status = ""
	StatusApplied Status = "applied"
	StatusRunning Status = "running"
	StatusTodo    Status = "todo"
	StatusFailed  Status = "failed"
	StatusUnknown Status = "unknown"
)

type SendStatus func(newStatus Status, info string, detail Detail, outVariables variables.Variables)
