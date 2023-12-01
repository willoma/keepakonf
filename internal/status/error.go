package status

import "io"

type Error struct {
	Err error
}

func (e Error) JSON(w io.StringWriter) {
	startDetailJSON(w, "error")

	w.WriteString(`"`)
	w.WriteString(e.Err.Error())
	w.WriteString(`"`)

	endDetailJSON(w)
}
