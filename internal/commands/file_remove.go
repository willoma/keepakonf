package commands

import (
	"fmt"
	"os"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = registerFileWatcher(
	"file remove",
	"remove",
	"Remove file or directory",
	ParamsDesc{
		{"path", "File path", ParamTypeFilePath},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCommand {
		return &fileRemove{
			fileWatcherCmdInit(params, vars, msg),
		}
	},
)

type fileRemove struct {
	fileWatcherCmd
}

func (f *fileRemove) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusDirectory, external.FileStatusFile:
		f.msg(status.StatusTodo, fmt.Sprintf("Need to remove %q", f.getPath()), nil, nil)
	case external.FileStatusUnknown:
		f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.getPath()), nil, nil)
	case external.FileStatusNotFound:
		f.msg(status.StatusApplied, fmt.Sprintf("%q does not exist", f.getPath()), nil, nil)
	}
}

func (f *fileRemove) apply() bool {
	if err := os.RemoveAll(f.getPath()); err != nil {
		f.msg(status.StatusFailed, fmt.Sprintf("Failed removing %q", f.getPath()), status.Error(err.Error()), nil)
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q removed", f.getPath()), nil, nil)
	return true
}
