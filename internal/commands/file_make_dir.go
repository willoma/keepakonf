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
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = registerFileWatcher(
	"file make dir",
	"folder",
	"Make directory",
	ParamsDesc{
		{"path", "Directory path", ParamTypeFilePath},
		{"owner", "Directory owner", ParamTypeUsername},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCommand {
		return &fileMakeDir{
			fileWatcherCmdInit(params, vars, msg),
			params["owner"].(string),
		}
	},
)

type fileMakeDir struct {
	fileWatcherCmd

	owner string
}

func (f *fileMakeDir) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusDirectory:
		// TODO check if the directory belongs to its owner
		f.msg(status.StatusApplied, fmt.Sprintf("%q exists", f.getPath()), nil, nil)
	case external.FileStatusFile:
		f.msg(status.StatusFailed, fmt.Sprintf("%q is not a directory", f.getPath()), nil, nil)
	case external.FileStatusUnknown:
		f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.getPath()), nil, nil)
	case external.FileStatusNotFound:
		f.msg(status.StatusTodo, fmt.Sprintf("Need to create %q", f.getPath()), nil, nil)
	}
}

func (f *fileMakeDir) apply() bool {
	finfo, err := os.Stat(f.getPath())
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", f.getPath()), status.Error(err.Error()), nil)
			return false
		}

		ownerUser, err := user.Lookup(f.vars.Replace(f.owner))
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not get user information for %q", f.owner), status.Error(err.Error()), nil)
			return false
		}

		if err := os.MkdirAll(f.getPath(), 0o755); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not create %q", f.getPath()), status.Error(err.Error()), nil)
			return false
		}

		uid, _ := strconv.Atoi(ownerUser.Uid)
		gid, _ := strconv.Atoi(ownerUser.Gid)
		if err := os.Chown(f.getPath(), uid, gid); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not change ownership of %q", f.getPath()), status.Error(err.Error()), nil)
			return false
		}

		f.msg(status.StatusApplied, fmt.Sprintf("%q directory created", f.getPath()), nil, nil)
		return true
	}

	if !finfo.IsDir() {
		f.msg(status.StatusFailed, fmt.Sprintf("%q is not a directory", f.getPath()), nil, nil)
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q exists", f.getPath()), nil, nil)
	return true
}
