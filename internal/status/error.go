package status

import "encoding/json"

type Error struct {
	Err error
}

func (e Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Err.Error())
}

func (e Error) DetailType() string {
	return "error"
}
