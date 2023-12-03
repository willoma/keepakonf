package status

import (
	"bytes"
	"encoding/json"
)

type DetailType string

type Detail interface {
	JSON() json.RawMessage
}

func startDetailJSON(detailType string) *bytes.Buffer {
	var w bytes.Buffer
	w.WriteString(`{"t":"`)
	w.WriteString(detailType)
	w.WriteString(`","d":`)
	return &w
}

func endDetailJSON(w *bytes.Buffer) json.RawMessage {
	w.WriteString("}")
	return json.RawMessage(w.Bytes())
}
