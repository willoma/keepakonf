package status

import (
	"io"
	"strings"
)

// Terminal JSON content:
//
// - cmd: Optional command line
// - out: Terminal output

type Terminal struct {
	Command string
	Output  string
}

func (t *Terminal) JSON(w io.StringWriter) {
	startDetailJSON(w, "terminal")

	w.WriteString(`{"`)

	if t.Command != "" {
		cmd := strings.ReplaceAll(t.Command, "\n", `\n`)
		cmd = strings.ReplaceAll(cmd, `"`, `\"`)
		w.WriteString(`cmd":"`)
		w.WriteString(cmd)
		w.WriteString(`","`)
	}

	out := strings.ReplaceAll(t.Output, "\n", `\n`)
	out = strings.ReplaceAll(out, `"`, `\"`)
	w.WriteString(`out": "`)
	w.WriteString(out)
	w.WriteString(`"}`)

	endDetailJSON(w)
}
