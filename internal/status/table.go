package status

import (
	"io"
	"strings"
)

// Table JSON content:
//
// - h: Optional headers (list of strings)
// - r: Rows (list of the following)
//   - obj.s: Cell status (cf status.Status)
//   - obj.c: Cell content

type TableCell struct {
	Status  Status
	Content string
}

type Table struct {
	Header []string
	Rows   [][]TableCell
}

func (t *Table) AppendRow(cells ...TableCell) {
	t.Rows = append(t.Rows, cells)
}

func (t *Table) JSON(w io.StringWriter) {
	startDetailJSON(w, "table")

	w.WriteString(`{"`)

	if len(t.Header) > 0 {
		w.WriteString(`h":["`)
		cell := strings.ReplaceAll(t.Header[0], "\n", `\n`)
		cell = strings.ReplaceAll(cell, `"`, `\"`)
		w.WriteString(cell)

		for _, h := range t.Header[1:] {
			w.WriteString(`","`)
			cell = strings.ReplaceAll(h, "\n", `\n`)
			cell = strings.ReplaceAll(cell, `"`, `\"`)
			w.WriteString(cell)
		}
		w.WriteString(`"],"r":[`)
	}

	if len(t.Rows) > 0 {
		w.WriteString("[")
		writeTableRow(w, t.Rows[0])

		for _, r := range t.Rows[1:] {
			w.WriteString("],[")
			writeTableRow(w, r)
		}

		w.WriteString("]")
	}

	w.WriteString("]}")
	endDetailJSON(w)
}

func writeTableRow(w io.StringWriter, r []TableCell) {
	if len(r) == 0 {
		return
	}

	w.WriteString(`{"s":"`)
	w.WriteString(string(r[0].Status))
	w.WriteString(`","c":"`)

	cell := strings.ReplaceAll(r[0].Content, "\n", `\n`)
	cell = strings.ReplaceAll(cell, `"`, `\"`)
	w.WriteString(cell)

	for _, c := range r[1:] {
		w.WriteString(`"},{"s":"`)
		w.WriteString(string(c.Status))
		w.WriteString(`","c":"`)

		cell := strings.ReplaceAll(c.Content, "\n", `\n`)
		cell = strings.ReplaceAll(cell, `"`, `\"`)
		w.WriteString(cell)
	}

	w.WriteString(`"}`)
}
