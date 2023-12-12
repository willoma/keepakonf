package status

type Error string

func (e Error) DetailType() string {
	return "error"
}
