package status

import (
	"bytes"
	"encoding/json"
)

type Detail interface {
	DetailType() string
}

func DetailJSON(d Detail) json.RawMessage {
	result, err := json.Marshal(struct {
		Type   string `json:"t"`
		Detail any    `json:"d"`
	}{
		Type:   d.DetailType(),
		Detail: d,
	})
	if err != nil {
		var b bytes.Buffer
		b.WriteString(`{"t":"`)
		b.WriteString(d.DetailType())
		b.WriteString(`","d":"`)
		b.WriteString(err.Error())
		b.WriteString(`"}`)
		return json.RawMessage(b.Bytes())
	}

	return json.RawMessage(result)
}
