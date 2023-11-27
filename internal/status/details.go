package status

import "io"

type DetailType string

type Detail interface {
	JSON(io.StringWriter)
}

func startDetailJSON(w io.StringWriter, detailType string) {
	w.WriteString(`{"t":"`)
	w.WriteString(detailType)
	w.WriteString(`","d":`)
}

func endDetailJSON(w io.StringWriter) {
	w.WriteString("}")
}
