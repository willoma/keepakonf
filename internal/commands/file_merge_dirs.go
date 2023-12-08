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
	func(params map[string]any, msg status.SendStatus) Command {
		return &fileMergeDirs{
			msg:         msg,
			source:      params["source"].(string),
			destination: params["destination"].(string),
			owner:       params["owner"].(string),

			closeChan: make(chan struct{}),
		}
	},
)

type fileMergeDirs struct {
	msg         status.SendStatus
	source      string
	destination string
	owner       string

	applying  atomic.Bool
	srcClose  func()
	dstClose  func()
	closeChan chan struct{}
}

func (f *fileMergeDirs) Watch() {
	srcChan, srcClose := external.WatchFile(f.source)
	f.srcClose = srcClose

	dstChan, dstClose := external.WatchFile(f.destination)
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
					Content: f.source,
				},
				status.TableCell{
					Status:  status.StatusTodo,
					Content: "To merge",
				},
			)
		case external.FileStatusFile, external.FileStatusFileChange:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusFailed,
					Content: f.source,
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
					Content: f.source,
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
					Content: f.source,
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
					Content: f.destination,
				},
				status.TableCell{
					Status:  status.StatusApplied,
					Content: "Exists",
				},
			)
		case external.FileStatusFile, external.FileStatusFileChange:
			result.AppendRow(
				status.TableCell{
					Status:  status.StatusFailed,
					Content: f.destination,
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
					Content: f.destination,
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
					Content: f.destination,
				},
				status.TableCell{
					Status:  status.StatusApplied,
					Content: "To create",
				},
			)
		}

		switch {
		case dstStatus == external.FileStatusFile, dstStatus == external.FileStatusFileChange:
			f.msg(status.StatusFailed, fmt.Sprintf("destination %q is not a directory", f.destination), &result)
		case srcStatus == external.FileStatusFile, srcStatus == external.FileStatusFileChange:
			f.msg(status.StatusFailed, fmt.Sprintf("source %q is not a directory", f.source), &result)
		case dstStatus == external.FileStatusUnknown:
			f.msg(status.StatusUnknown, fmt.Sprintf("destination %q status unknown", f.destination), &result)
		case srcStatus == external.FileStatusUnknown:
			f.msg(status.StatusUnknown, fmt.Sprintf("source %q status unknown", f.source), &result)
		case srcStatus == external.FileStatusDirectory && dstStatus == external.FileStatusDirectory:
			f.msg(status.StatusTodo, fmt.Sprintf("Need to merge source %q into destination %q", f.source, f.destination), &result)
		case srcStatus == external.FileStatusDirectory && dstStatus == external.FileStatusNotFound:
			f.msg(status.StatusTodo, fmt.Sprintf("Need to move source %q to destination %q", f.source, f.destination), &result)
		case srcStatus == external.FileStatusNotFound && dstStatus == external.FileStatusNotFound:
			f.msg(status.StatusTodo, fmt.Sprintf("Need to create %q", f.destination), &result)
		case srcStatus == external.FileStatusNotFound && dstStatus == external.FileStatusDirectory:
			f.msg(status.StatusApplied, fmt.Sprintf("destination %q exists", f.destination), &result)
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

	if finfo, err := os.Stat(f.source); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Source does not exist, simply make sure destination exists
			return f.confirmDestination()
		} else {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of source %q", f.source), status.Error{Err: err})
			return false
		}
	} else if !finfo.IsDir() {
		f.msg(status.StatusFailed, fmt.Sprintf("Source %q is not a directory", f.source), nil)
		return false
	}

	// Now we know source exists and is a directory
	if !f.mergeDir(f.source, f.destination) {
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q merged into destination %q", f.source, f.destination), nil)
	return true
}

// confirmDestination is called only when source does not exist, to confirm destination exists
func (f *fileMergeDirs) confirmDestination() bool {
	finfo, err := os.Stat(f.destination)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", f.destination), status.Error{Err: err})
			return false
		}

		ownerUser, err := user.Lookup(f.owner)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not get user information for %q", f.owner), status.Error{Err: err})
			return false
		}

		if err := os.MkdirAll(f.destination, 0o755); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not create destination %q", f.destination), status.Error{Err: err})
			return false
		}

		uid, _ := strconv.Atoi(ownerUser.Uid)
		gid, _ := strconv.Atoi(ownerUser.Gid)
		if err := os.Chown(f.destination, uid, gid); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not change ownership of destination %q", f.destination), status.Error{Err: err})
			return false
		}

		f.msg(status.StatusApplied, fmt.Sprintf("%q directory created, and %q does not exist, nothing to merge", f.destination, f.source), nil)
		return true
	}

	if !finfo.IsDir() {
		f.msg(status.StatusFailed, fmt.Sprintf("Destination %q is not a directory", f.destination), nil)
		return false
	}

	f.msg(status.StatusApplied, fmt.Sprintf("%q does not exist, nothing to merge into %q", f.source, f.destination), nil)
	return true
}

// mergeDir merges source into destination, recursively
func (f *fileMergeDirs) mergeDir(srcdir, dstdir string) bool {
	entries, err := os.ReadDir(srcdir)
	if err != nil {
		f.msg(status.StatusFailed, fmt.Sprintf("Could not read content of %q", srcdir), status.Error{Err: err})
		return false
	}

	for _, e := range entries {
		fname := e.Name()

		srcPath := filepath.Join(srcdir, fname)
		dstPath := filepath.Join(dstdir, fname)

		dstFinfo, err := os.Stat(dstPath)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				f.msg(status.StatusFailed, fmt.Sprintf("Could not check status of %q", dstPath), status.Error{Err: err})
				return false
			}

			// Destination does not exist, simply move source
			if err := os.Rename(srcPath, dstPath); err != nil {
				f.msg(status.StatusFailed, fmt.Sprintf("Could not move %q to %q", srcPath, dstdir), status.Error{Err: err})
				return false
			}
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
			f.msg(status.StatusFailed, fmt.Sprintf("Could not move directory %q to %q, which is not a directory", srcPath, dstPath), nil)
			return false
		case dstFinfo.IsDir():
			f.msg(status.StatusFailed, fmt.Sprintf("Could not move file %q to %q, which is a directory", srcPath, dstPath), nil)
			return false
		}

		// From here, we know both are NOT directories, problems may arise :-)
		srcData, err := os.ReadFile(srcPath)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not read content of source %q", srcPath), status.Error{Err: err})
			return false
		}
		srcSum := sha256.Sum256(srcData)

		dstData, err := os.ReadFile(dstPath)
		if err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not read content of destination %q", dstPath), status.Error{Err: err})
			return false
		}
		dstSum := sha256.Sum256(dstData)

		if srcSum != dstSum {
			// Files are different, we do not know what to do
			f.msg(status.StatusFailed, fmt.Sprintf("Source %q and destination %q are different", srcPath, dstPath), status.Error{Err: err})
			return false
		}

		if err := os.Remove(srcPath); err != nil {
			f.msg(status.StatusFailed, fmt.Sprintf("Could not remove source %q", srcPath), status.Error{Err: err})
			return false
		}
	}

	return true
}
