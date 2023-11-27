package status

type Status string

const (
	StatusFailed  Status = "failed"
	StatusNone    Status = "none"
	StatusTodo    Status = "todo"
	StatusRunning Status = "running"
	StatusApplied Status = "applied"
	StatusUnknown Status = "unknown"
)

var priority = map[Status]int{
	StatusFailed:  5,
	StatusNone:    4,
	StatusTodo:    3,
	StatusRunning: 2,
	StatusApplied: 1,
	StatusUnknown: 0,
}

func (s Status) IfHigherPriority(newS Status) Status {
	if priority[s] > priority[newS] {
		return s
	}
	return newS
}

type SendStatus func(newStatus Status, info string, detail Detail)
