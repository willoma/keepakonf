package status

type Terminal struct {
	Command string `json:"cmd,omitempty"`
	Output  string `json:"out"`
}

func (t *Terminal) DetailType() string {
	return "terminal"
}
