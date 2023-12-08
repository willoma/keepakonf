package commands

import (
	"fmt"
	"os"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
)

var _ = registerFileWatcher(
	"file remove",
	"remove",
	"Remove file or directory",
	ParamsDesc{
		{"path", "File path", ParamTypeFilePath},
	},
	func(params map[string]any, msg status.SendStatus) fileWatcherCommand {
		return &fileRemove{
			path: params["path"].(string),
			msg:  msg,
		}
	},
)

type fileRemove struct {
	msg  status.SendStatus
	path string
}

func (f *fileRemove) getPath() string {
	return f.path
}

func (f *fileRemove) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusDirectory, external.FileStatusFile, external.FileStatusFileChange:
		f.msg(status.StatusTodo, fmt.Sprintf("Need to remove %q", f.path), nil)
	case external.FileStatusUnknown:
		f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.path), nil)
	case external.FileStatusNotFound:
		f.msg(status.StatusApplied, fmt.Sprintf("%q does not exist", f.path), nil)
	}
}

func (f *fileRemove) apply() bool {
	if err := os.RemoveAll(f.path); err != nil {
		f.msg(status.StatusFailed, fmt.Sprintf("Failed removing %q", f.path), status.Error{Err: err})
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q removed", f.path), nil)
	return true
}
