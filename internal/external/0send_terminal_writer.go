package external

import (
	"strings"

	"github.com/willoma/keepakonf/internal/status"
)

type sendTerminalWriter struct {
	cmdline    string
	lineBuffer []byte
	result     strings.Builder
	receiver   func(status.Status, string, status.Detail)
}

func newSendWriter(receiver func(status.Status, string, status.Detail), cmdline string) *sendTerminalWriter {
	return &sendTerminalWriter{
		cmdline:    cmdline,
		lineBuffer: make([]byte, 0, 80),
		receiver:   receiver,
	}
}

func (w *sendTerminalWriter) Write(p []byte) (n int, err error) {
	cr := false
	for _, b := range p {
		switch b {
		case '\n':
			if cr {
				cr = false
			}
			w.lineBuffer = append(w.lineBuffer, '\n')
			if c, err := w.result.Write(w.lineBuffer); err != nil {
				return c, err
			}
			w.lineBuffer = w.lineBuffer[:0]

			w.receiver(
				status.StatusRunning,
				"",
				&status.Terminal{
					Command: w.cmdline,
					Output:  w.result.String(),
				},
			)
		case '\r':
			cr = true
		default:
			if cr {
				w.lineBuffer = w.lineBuffer[:0]
				cr = false
			}
			w.lineBuffer = append(w.lineBuffer, b)
		}
	}
	return len(p), nil
}

func (w *sendTerminalWriter) Result() string {
	return w.result.String()
}
