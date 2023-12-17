package log

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/willoma/keepakonf/internal/status"
)

type messageStruct struct {
	Timestamp     string          `json:"ts"`
	Message       string          `json:"msg"`
	Icon          string          `json:"ico,omitempty"`
	Status        status.Status   `json:"st,omitempty"`
	GroupID       string          `json:"gid,omitempty"`
	InstructionID string          `json:"iid,omitempty"`
	GroupName     string          `json:"grp,omitempty"`
	Detail        json.RawMessage `json:"dtl,omitempty"`
}

func write(
	msg string,
	icon string,
	status status.Status,
	groupID string,
	instructionID string,
	groupName string,
	detail json.RawMessage,
) {
	logMsg := messageStruct{
		Timestamp:     time.Now().Format("2006-01-02T15:04:05"),
		Message:       msg,
		Icon:          icon,
		Status:        status,
		GroupID:       groupID,
		InstructionID: instructionID,
		GroupName:     groupName,
		Detail:        detail,
	}

	jsonRec, err := json.Marshal(logMsg)
	if err != nil {
		slog.Error("Could not convert log to JSON", "log", logMsg, "error", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		slog.Error("Could not open log file", "error", err)
	} else {
		if _, err := f.Write(append(jsonRec, '\n')); err != nil {
			slog.Error("could not write log to file", "error", err)
		}
		f.Close()
	}

	records = append(records, jsonRec)
	if io != nil {
		io.Emit("log", json.RawMessage(jsonRec))
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
	write(
		msg,
		"error",
		status.StatusFailed,
		"", "", "",
		status.DetailJSON(
			status.Error(err.Error()),
		),
	)
}

func Errorf(err error, format string, args ...any) {
	write(
		fmt.Sprintf(format, args...),
		"error",
		status.StatusFailed,
		"", "", "",
		status.DetailJSON(
			status.Error(err.Error()),
		),
	)
}
