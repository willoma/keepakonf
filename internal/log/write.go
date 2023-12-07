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

func write(
	msg string,
	icon string,
	status status.Status,
	groupID string,
	instructionID string,
	groupName string,
	detail json.RawMessage,
) {
	jsonRec := rawJSON(msg, icon, status, groupID, instructionID, groupName, detail)

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		Error(err, "Could not open log file")
		slog.Error("Could not open log file", "error", err)
	} else {
		if _, err := f.Write(append(jsonRec, '\n')); err != nil {
			Error(err, "Could not write to log file")
			slog.Error("could not write log to file", "error", err)
		}
		f.Close()
	}

	records = append(records, jsonRec)
	if io != nil {
		io.Emit("log", jsonRec)
	}
}

func Info(
	msg string,
	icon string,
	status status.Status,
	groupID string,
	instructionID string,
	groupName string,
	detail json.RawMessage,
) {
	write(msg, icon, status, groupID, instructionID, groupName, detail)
}

func Error(err error, msg string) {
	write(msg, "error", status.StatusFailed, "", "", "", status.Error{Err: err}.JSON())
}

func Errorf(err error, format string, args ...any) {
	write(fmt.Sprintf(format, args...), "error", status.StatusFailed, "", "", "", status.Error{Err: err}.JSON())
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
	detail json.RawMessage,
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
		raw.Write(detail)
	}

	raw.WriteByte('}')

	return raw.Bytes()
}
