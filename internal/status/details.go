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

type Error string

func (e Error) DetailType() string {
	return "error"
}

type TableCell struct {
	Status  Status `json:"s"`
	Content string `json:"c"`
}

type Table struct {
	Header []string      `json:"h,omitempty"`
	Rows   [][]TableCell `json:"r"`
}

func (t *Table) DetailType() string {
	return "table"
}

func (t *Table) AppendRow(cells ...TableCell) {
	t.Rows = append(t.Rows, cells)
}

type Terminal struct {
	Command string `json:"cmd,omitempty"`
	Output  string `json:"out"`
}

func (t *Terminal) DetailType() string {
	return "terminal"
}
