package commands

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = register(
	"file merge dirs",
	"folder",
	"Merge directories",
	ParamsDesc{
		{"source", "Source directory", ParamTypeFilePath},
		{"destination", "Destination directory", ParamTypeFilePath},
		{"owner", "Destination dir owner", ParamTypeUsername},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) Command {
		return &fileMergeDirs{
			msg:         msg,
			vars:        vars,
			source:      params["source"].(string),
			destination: params["destination"].(string),
			owner:       params["owner"].(string),

			closeChan: make(chan struct{}),
		}
	},
)

type fileMergeDirs struct {
	msg  status.SendStatus
	vars variables.Variables

	source      string
	destination string
	owner       string

	applying  atomic.Bool
	srcClose  func()
	dstClose  func()
	closeChan chan struct{}
}

func (f *fileMergeDirs) UpdateVariables(vars variables.Variables) {
	if f.vars.Update(vars) {
		f.Stop()
		f.Watch()
	}
}

func (f *fileMergeDirs) Watch() {
	srcChan, srcClose := external.WatchFile(f.vars.Replace(f.source))
	f.srcClose = srcClose

	dstChan, dstClose := external.WatchFile(f.vars.Replace(f.destination))
	f.dstClose = dstClose

	var srcStatus, dstStatus external.FileStatus

	check := func() {
		if f.applying.Load() {
			// No update if it is currently applying
			return
		}

		result := status.Table{
			Header: []string{"Directory", "Status"},
		}

		switch srcStatus {
		case external.FileStatusDirectory:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusTodo,
					Content: f.vars.Replace(f.source),
				},
				status.TableCell{
					Status:  status.StatusTodo,
					Content: "To merge",
				},
			)
		case external.FileStatusFile:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusFailed,
					Content: f.vars.Replace(f.source),
				},
				status.TableCell{
					Status:  status.StatusFailed,
					Content: "Not a directory",
				},
			)
		case external.FileStatusUnknown:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusUnknown,
					Content: f.vars.Replace(f.source),
				},
				status.TableCell{
					Status:  status.StatusUnknown,
					Content: "Unknown",
				},
			)
		case external.FileStatusNotFound:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusApplied,
					Content: f.vars.Replace(f.source),
				},
				status.TableCell{
					Status:  status.StatusApplied,
					Content: "Does not exist",
				},
			)
		}

		switch dstStatus {
		case external.FileStatusDirectory:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusApplied,
					Content: f.vars.Replace(f.destination),
				},
				status.TableCell{
					Status:  status.StatusApplied,
					Content: "Exists",
				},
			)
			// TODO check if the directory belongs to its owner
		case external.FileStatusFile:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusFailed,
					Content: f.vars.Replace(f.destination),
				},
				status.TableCell{
					Status:  status.StatusFailed,
					Content: "Not a directory",
				},
			)
		case external.FileStatusUnknown:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusUnknown,
					Content: f.vars.Replace(f.destination),
				},
				status.TableCell{
					Status:  status.StatusUnknown,
					Content: "Unknown",
				},
			)
		case external.FileStatusNotFound:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusApplied,
					Content: f.vars.Replace(f.destination),
				},
				status.TableCell{
					Status:  status.StatusApplied,
					Content: "To create",
				},
			)
		}

		switch {
		case dstStatus == external.FileStatusFile:
			f.msg(status.StatusFailed, fmt.Sprintf("destination %q is not a directory", f.vars.Replace(f.destination)), &result, nil)
		case srcStatus == external.FileStatusFile:
			f.msg(status.StatusFailed, fmt.Sprintf("source %q is not a directory", f.vars.Replace(f.source)), &result, nil)
		case dstStatus == external.FileStatusUnknown:
			f.msg(status.StatusUnknown, fmt.Sprintf("destination %q status unknown", f.vars.Replace(f.destination)), &result, nil)
		case srcStatus == external.FileStatusUnknown:
			f.msg(status.StatusUnknown, fmt.Sprintf("source %q status unknown", f.vars.Replace(f.source)), &result, nil)
		case srcStatus == external.FileStatusDirectory && dstStatus == external.FileStatusDirectory:
			f.msg(status.StatusTodo, fmt.Sprintf("Need to merge source %q into destination %q", f.vars.Replace(f.source), f.vars.Replace(f.destination)), &result, nil)
		case srcStatus == external.FileStatusDirectory && dstStatus == external.FileStatusNotFound:
			f.msg(status.StatusTodo, fmt.Sprintf("Need to move source %q to destination %q", f.vars.Replace(f.source), f.vars.Replace(f.destination)), &result, nil)
		case srcStatus == external.FileStatusNotFound && dstStatus == external.FileStatusNotFound:
			f.msg(status.StatusTodo, fmt.Sprintf("Need to create %q", f.vars.Replace(f.destination)), &result, nil)
		case srcStatus == external.FileStatusNotFound && dstStatus == external.FileStatusDirectory:
			f.msg(status.StatusApplied, fmt.Sprintf("destination %q exists", f.vars.Replace(f.destination)), &result, nil)
		}
	}

	go func() {
		for {
			select {
			case srcStatus = <-srcChan:
				check()
			case dstStatus = <-dstChan:
				check()
			case <-f.closeChan:
				return
			}
		}
	}()
}

func (f *fileMergeDirs) Stop() {
	if f.srcClose != nil {
		f.srcClose()
	}
	if f.dstClose != nil {
		f.dstClose()
	}
	f.closeChan <- struct{}{}
}

func (f *fileMergeDirs) Apply() bool {
	f.applying.Store(true)
	defer f.applying.Store(false)

	if finfo, err := os.Stat(f.vars.Replace(f.source)); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Source does not exist, simply make sure destination exists
			return f.confirmDestination()
		} else {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of source %q", f.vars.Replace(f.source)), status.Error(err.Error()), nil)
			return false
		}
	} else if !finfo.IsDir() {
		f.msg(status.StatusFailed, fmt.Sprintf("Source %q is not a directory", f.vars.Replace(f.source)), nil, nil)
		return false
	}

	// Now we know source exists and is a directory
	if !f.mergeDir(f.vars.Replace(f.source), f.vars.Replace(f.destination)) {
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q merged into destination %q", f.vars.Replace(f.source), f.vars.Replace(f.destination)), nil, nil)
	return true
}

// confirmDestination is called only when source does not exist, to confirm destination exists
func (f *fileMergeDirs) confirmDestination() bool {
	finfo, err := os.Stat(f.vars.Replace(f.destination))
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", f.vars.Replace(f.destination)), status.Error(err.Error()), nil)
			return false
		}

		ownerUser, err := user.Lookup(f.vars.Replace(f.owner))
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not get user information for %q", f.owner), status.Error(err.Error()), nil)
			return false
		}

		if err := os.MkdirAll(f.vars.Replace(f.destination), 0o755); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not create destination %q", f.vars.Replace(f.destination)), status.Error(err.Error()), nil)
			return false
		}

		uid, _ := strconv.Atoi(ownerUser.Uid)
		gid, _ := strconv.Atoi(ownerUser.Gid)
		if err := os.Chown(f.vars.Replace(f.destination), uid, gid); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not change ownership of destination %q", f.vars.Replace(f.destination)), status.Error(err.Error()), nil)
			return false
		}

		f.msg(status.StatusApplied, fmt.Sprintf("%q directory created, and %q does not exist, nothing to merge", f.vars.Replace(f.destination), f.vars.Replace(f.source)), nil, nil)
		return true
	}

	if !finfo.IsDir() {
		f.msg(status.StatusFailed, fmt.Sprintf("Destination %q is not a directory", f.vars.Replace(f.destination)), nil, nil)
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q does not exist, nothing to merge into %q", f.vars.Replace(f.source), f.vars.Replace(f.destination)), nil, nil)
	return true
}

// mergeDir merges source into destination, recursively
func (f *fileMergeDirs) mergeDir(srcdir, dstdir string) bool {
	entries, err := os.ReadDir(srcdir)
	if err != nil {
		f.msg(status.StatusFailed, fmt.Sprintf("Could not read content of %q", srcdir), status.Error(err.Error()), nil)
		return false
	}

	for _, e := range entries {
		fname := e.Name()

		srcPath := filepath.Join(srcdir, fname)
		dstPath := filepath.Join(dstdir, fname)

		dstFinfo, err := os.Stat(dstPath)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", dstPath), status.Error(err.Error()), nil)
				return false
			}

			// Destination does not exist, simply move source
			if err := os.Rename(srcPath, dstPath); err != nil {
				f.msg(status.StatusFailed, fmt.Sprintf("Could not move %q to %q", srcPath, dstdir), status.Error(err.Error()), nil)
				return false
			}
			// TODO change ownership if needed (only for root of merged dir)
			return true
		}

		// From here, we know destination exists, we must check and merge

		switch {
		case e.IsDir() && dstFinfo.IsDir():
			// Both are directories, merge their contents
			if !f.mergeDir(srcPath, dstPath) {
				return false
			}
		case e.IsDir():
			// Source is directory, but destination is not
			f.msg(status.StatusFailed, fmt.Sprintf("Could not move directory %q to %q, which is not a directory", srcPath, dstPath), nil, nil)
			return false
		case dstFinfo.IsDir():
			f.msg(status.StatusFailed, fmt.Sprintf("Could not move file %q to %q, which is a directory", srcPath, dstPath), nil, nil)
			return false
		}

		// From here, we know both are NOT directories, problems may arise :-)
		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not read content of source %q", srcPath), status.Error(err.Error()), nil)
			return false
		}
		srcSum := sha256.Sum256(srcData)

		dstData, err := os.ReadFile(dstPath)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not read content of destination %q", dstPath), status.Error(err.Error()), nil)
			return false
		}
		dstSum := sha256.Sum256(dstData)

		if srcSum != dstSum {
			// Files are different, we do not know what to do
			f.msg(status.StatusFailed, fmt.Sprintf("Source %q and destination %q are different", srcPath, dstPath), status.Error(err.Error()), nil)
			return false
		}

		if err := os.Remove(srcPath); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not remove source %q", srcPath), status.Error(err.Error()), nil)
			return false
		}
	}

	return true
}
