package commands

import (
	"fmt"
	"os"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = registerFileWatcher(
	"file content",
	"file",
	"Ensure a file has a content",
	ParamsDesc{
		{"path", "File path", ParamTypeFilePath},
		{"content", "File content", ParamTypeText},
		{"owner", "File owner", ParamTypeUsername},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCommand {
		return &fileContent{
			fileWatcherCmdInit(params, vars, msg),
			params["content"].(string),
			params["owner"].(string),
		}
	},
)

type fileContent struct {
	fileWatcherCmd

	content string
	owner   string
}

func (f *fileContent) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusDirectory:
		f.msg(status.StatusFailed, fmt.Sprintf("%q is a directory", f.getPath()), nil, nil)
	case external.FileStatusFile:
		currentB, err := os.ReadFile(f.getPath())
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("could not read %q", f.getPath()), status.Error(err.Error()), nil)
			return
		}
		current := string(currentB)
		if current != f.vars.Replace(f.content) {
			f.msg(status.StatusTodo, fmt.Sprintf("Need to change %q", f.getPath()), status.TextDiff{Before: current, After: f.vars.Replace(f.content)}, nil)
			return
		}
		// TODO check ownership
		f.msg(status.StatusApplied, fmt.Sprintf("%q has the required content", f.getPath()), status.Text(current), nil)
	case external.FileStatusUnknown:
		f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.getPath()), nil, nil)
	case external.FileStatusNotFound:
		f.msg(status.StatusTodo, fmt.Sprintf("Need to create %q", f.getPath()), status.Text(f.vars.Replace(f.content)), nil)
	}
}

func (f *fileContent) apply() bool {
	if err := os.WriteFile(f.getPath(), []byte(f.vars.Replace(f.content)), 0x644); err != nil {
		f.msg(status.StatusFailed, fmt.Sprintf("Could not write to %q", f.getPath()), status.Error(err.Error()), nil)
		return false
	}
	if f.owner != "" {
		userData, err := external.GetUser(f.vars.Replace(f.owner))
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not get user information for %q", f.vars.Replace(f.owner)), status.Error(err.Error()), nil)
			return false
		}
		if err := os.Chown(f.getPath(), userData.ID, userData.GID); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not change %q ownership to %q", f.getPath(), f.vars.Replace(f.owner)), status.Error(err.Error()), nil)
			return false
		}
	}

	f.msg(status.StatusApplied, fmt.Sprintf("Wrote content to %q", f.getPath()), status.Text(f.content), nil)
	return true
}
