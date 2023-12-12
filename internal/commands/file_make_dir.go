package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
)

var _ = registerFileWatcher(
	"file make dir",
	"folder",
	"Make directory",
	ParamsDesc{
		{"path", "Directory path", ParamTypeFilePath},
		{"owner", "Directory owner", ParamTypeUsername},
	},
	func(params map[string]any, msg status.SendStatus) fileWatcherCommand {
		return &fileMakeDir{
			msg:   msg,
			path:  params["path"].(string),
			owner: params["owner"].(string),
		}
	},
)

type fileMakeDir struct {
	msg   status.SendStatus
	path  string
	owner string
}

func (f *fileMakeDir) getPath() string {
	return f.path
}

func (f *fileMakeDir) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusDirectory:
		f.msg(status.StatusApplied, fmt.Sprintf("%q exists", f.path), nil)
	case external.FileStatusFile:
		f.msg(status.StatusFailed, fmt.Sprintf("%q is not a directory", f.path), nil)
	case external.FileStatusUnknown:
		f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.path), nil)
	case external.FileStatusNotFound:
		f.msg(status.StatusTodo, fmt.Sprintf("Need to create %q", f.path), nil)
	}
}

func (f *fileMakeDir) apply() bool {
	finfo, err := os.Stat(f.path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", f.path), status.Error(err.Error()))
			return false
		}

		ownerUser, err := user.Lookup(f.owner)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not get user information for %q", f.owner), status.Error(err.Error()))
			return false
		}

		if err := os.MkdirAll(f.path, 0o755); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not create %q", f.path), status.Error(err.Error()))
			return false
		}

		uid, _ := strconv.Atoi(ownerUser.Uid)
		gid, _ := strconv.Atoi(ownerUser.Gid)
		if err := os.Chown(f.path, uid, gid); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not change ownership of %q", f.path), status.Error(err.Error()))
			return false
		}

		f.msg(status.StatusApplied, fmt.Sprintf("%q directory created", f.path), nil)
		return true
	}

	if !finfo.IsDir() {
		f.msg(status.StatusFailed, fmt.Sprintf("%q is not a directory", f.path), nil)
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q exists", f.path), nil)
	return true
}
