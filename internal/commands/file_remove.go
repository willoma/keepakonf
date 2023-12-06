package commands

import (
	"fmt"
	"os"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
)

var _ = register(
	"file remove",
	"folder",
	"Remove file or directory",
	ParamsDesc{
		{"path", "File path", ParamTypeFilePath},
	},
	func(params map[string]any, logger *log.Logger, msg status.SendStatus) Command {
		return &fileRemove{
			command: command{logger, msg},
			path:    params["path"].(string),
		}
	},
)

type fileRemove struct {
	command
	path string

	applying atomic.Bool
	close    func()
}

func (f *fileRemove) Watch() {
	signals, close := external.WatchFile(f.logger, f.path)
	f.close = close

	go func() {
		for fstatus := range signals {

			if f.applying.Load() {
				// No update if it is currently applying
				continue
			}

			switch fstatus {
			case external.FileStatusDirectory, external.FileStatusFile, external.FileStatusFileChange:
				f.msg(status.StatusTodo, fmt.Sprintf("Need to remove %q", f.path), nil)
			case external.FileStatusUnknown:
				f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.path), nil)
			case external.FileStatusNotFound:
				f.msg(status.StatusApplied, fmt.Sprintf("%q does not exist", f.path), nil)
			}
		}
	}()
}

func (f *fileRemove) Stop() {
	if f.close != nil {
		f.close()
	}
}

func (f *fileRemove) Apply() bool {
	f.applying.Store(true)
	defer f.applying.Store(false)

	if err := os.RemoveAll(f.path); err != nil {
		f.msg(status.StatusFailed, fmt.Sprintf("Failed removing %q", f.path), status.Error{Err: err})
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q removed", f.path), nil)
	return true
}
