package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"strconv"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
)

var _ = register(
	"file make dir",
	"folder",
	"Make directory",
	ParamsDesc{
		{"path", "Directory path", ParamTypeFilePath},
		{"owner", "Directory owner", ParamTypeUsername},
	},
	func(params map[string]any, msg status.SendStatus) Command {
		return &fileMakeDir{
			command: command{msg},
			path:    params["path"].(string),
			owner:   params["owner"].(string),
		}
	},
)

type fileMakeDir struct {
	command
	path  string
	owner string

	applying atomic.Bool
	close    func()
}

func (f *fileMakeDir) Watch() {
	signals, close := external.WatchFile(f.path)
	f.close = close

	go func() {
		for fstatus := range signals {

			if f.applying.Load() {
				// No update if it is currently applying
				continue
			}

			switch fstatus {
			case external.FileStatusDirectory:
				f.msg(status.StatusApplied, fmt.Sprintf("%q exists", f.path), nil)
			case external.FileStatusFile, external.FileStatusFileChange:
				f.msg(status.StatusFailed, fmt.Sprintf("%q is not a directory", f.path), nil)
			case external.FileStatusUnknown:
				f.msg(status.StatusUnknown, fmt.Sprintf("%q status unknown", f.path), nil)
			case external.FileStatusNotFound:
				f.msg(status.StatusTodo, fmt.Sprintf("Need to create %q", f.path), nil)
			}
		}
	}()
}

func (f *fileMakeDir) Stop() {
	if f.close != nil {
		f.close()
	}
}

func (f *fileMakeDir) Apply() bool {
	f.applying.Store(true)
	defer f.applying.Store(false)

	finfo, err := os.Stat(f.path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", f.path), status.Error{Err: err})
			return false
		}

		ownerUser, err := user.Lookup(f.owner)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not get user information for %q", f.owner), status.Error{Err: err})
			return false
		}

		if err := os.MkdirAll(f.path, 0o755); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not create %q", f.path), status.Error{Err: err})
			return false
		}

		uid, _ := strconv.Atoi(ownerUser.Uid)
		gid, _ := strconv.Atoi(ownerUser.Gid)
		if err := os.Chown(f.path, uid, gid); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not change ownership of %q", f.path), status.Error{Err: err})
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
