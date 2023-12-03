package status

import (
	"encoding/json"
)

type Error struct {
	Err error
}

func (e Error) JSON() json.RawMessage {
	w := startDetailJSON("error")

	w.WriteString(`"`)
	w.WriteString(e.Err.Error())
	w.WriteString(`"`)

	return endDetailJSON(w)
}
