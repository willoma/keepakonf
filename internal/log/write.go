package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/willoma/keepakonf/internal/status"
)

// Log record JSON content:
//
// - ts:  Timestamp ("YYYY-MM-DDThh:mm:ss")
// - msg: Message
// - ico: Optional icon name
// - st:  Optional status (cf status.Status)
// - gid: Optional group ID
// - iid: Optional instruction ID
// - grp: Optional group name
// - dtl: Optional detail:
//   - dtl.t: Detail type name
//   - dtl.d: Detail content (different according to which status.Detail is provided)
// - err: Optional error

func (l *Logger) write(
	msg string,
	icon string,
	status status.Status,
	groupID string,
	instructionID string,
	groupName string,
	detail status.Detail,
) {
	jsonRec := rawJSON(msg, icon, status, groupID, instructionID, groupName, detail)

	l.mu.Lock()
	defer l.mu.Unlock()

	f, err := os.OpenFile(l.filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		l.Error(err, "Could not open log file")
		slog.Error("Could not open log file", "error", err)
	} else {
		if _, err := f.Write(append(jsonRec, '\n')); err != nil {
			l.Error(err, "Could not write to log file")
			slog.Error("could not write log to file", "error", err)
		}
		f.Close()
	}

	l.records = append(l.records, jsonRec)
	l.io.Emit("log", jsonRec)
}

func (l *Logger) Info(
	msg string,
	icon string,
	status status.Status,
	groupID string,
	instructionID string,
	groupName string,
	detail status.Detail,
) {
	l.write(msg, icon, status, groupID, instructionID, groupName, detail)
}

func (l *Logger) Error(err error, msg string) {
	l.write(msg, "error", status.StatusFailed, "", "", "", status.Error{Err: err})
}

func (l *Logger) Errorf(err error, format string, args ...any) {
	l.write(fmt.Sprintf(format, args...), "error", status.StatusFailed, "", "", "", status.Error{Err: err})
}

// rawJSON returns the record in JSON format. We do not use json.Marshal
// in order to improve performance and because we know exactly what type
// of data we can get.
//
// This performance improvement would probably not be visible to the end-user
// (except maybe with millions of log records), but I had fun writing this :-)
func rawJSON(
	msg string,
	icon string,
	status status.Status,
	groupID string,
	instructionID string,
	groupName string,
	detail status.Detail,
) json.RawMessage {
	var raw bytes.Buffer

	raw.WriteString(`{"ts":"`)
	raw.WriteString(time.Now().Format("2006-01-02T15:04:05"))

	msg = strings.ReplaceAll(msg, "\n", `\n`)
	msg = strings.ReplaceAll(msg, `"`, `\"`)
	raw.WriteString(`","msg":"`)
	raw.WriteString(msg)

	if icon != "" {
		raw.WriteString(`","ico":"`)
		raw.WriteString(icon)
	}

	if status != "" {
		raw.WriteString(`","st":"`)
		raw.WriteString(string(status))
	}

	if groupID != "" {
		raw.WriteString(`","gid":"`)
		raw.WriteString(groupID)
	}

	if instructionID != "" {
		raw.WriteString(`","iid":"`)
		raw.WriteString(instructionID)
	}

	if groupName != "" {
		raw.WriteString(`","grp":"`)
		raw.WriteString(groupName)
	}

	raw.WriteByte('"')

	if detail != nil {
		raw.WriteString(`,"dtl":`)
		detail.JSON(&raw)
	}

	raw.WriteByte('}')

	return raw.Bytes()
}
