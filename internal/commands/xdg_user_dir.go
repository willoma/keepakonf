package commands

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
)

const (
	xdgUserDirPath = ".config/user-dirs.dirs"
)

var (
	xdgUserDirRe   = regexp.MustCompile(`^XDG_(.*)_DIR="(.*)"`)
	xdgOrderedVars = []string{
		"DESKTOP",
		"DOWNLOAD",
		"TEMPLATES",
		"PUBLICSHARE",
		"DOCUMENTS",
		"MUSIC",
		"PICTURES",
		"VIDEOS",
	}
)

var _ = registerFileWatcher(
	"xdg user dir",
	"folder",
	"Set XDG user directories",
	ParamsDesc{
		{"user", "User", ParamTypeUsername},
		{"desktop", "Desktop", ParamTypeFilePath},
		{"download", "Download", ParamTypeFilePath},
		{"templates", "Templates", ParamTypeFilePath},
		{"publicshare", "Public share", ParamTypeFilePath},
		{"documents", "Documents", ParamTypeFilePath},
		{"music", "Music", ParamTypeFilePath},
		{"pictures", "Pictures", ParamTypeFilePath},
		{"videos", "Videos", ParamTypeFilePath},
	},
	func(params map[string]any, msg status.SendStatus) fileWatcherCommand {
		user := params["user"].(string)
		userData, err := external.GetUser(user)
		if err != nil {
			msg(status.StatusFailed, "Could not get user information for "+user, status.Error{Err: err})
		}
		return &xdgUserDir{
			msg:   msg,
			user:  userData,
			fpath: filepath.Join(userData.Home, xdgUserDirPath),
			required: map[string]string{
				"DESKTOP":     params["desktop"].(string),
				"DOWNLOAD":    params["download"].(string),
				"TEMPLATES":   params["templates"].(string),
				"PUBLICSHARE": params["publicshare"].(string),
				"DOCUMENTS":   params["documents"].(string),
				"MUSIC":       params["music"].(string),
				"PICTURES":    params["pictures"].(string),
				"VIDEOS":      params["videos"].(string),
			},
		}
	},
)

type xdgUserDir struct {
	msg status.SendStatus

	user  external.User
	fpath string

	required map[string]string
}

func (x *xdgUserDir) getPath() string {
	return x.fpath
}

func (x *xdgUserDir) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusFile, external.FileStatusFileChange:
		x.msg(x.check())
	case external.FileStatusDirectory:
		x.msg(status.StatusFailed, `"`+x.fpath+`" is a directory`, nil)
	case external.FileStatusNotFound:
		x.msg(status.StatusFailed, `File "`+x.fpath+`" not found`, nil)
	}
}

func (x *xdgUserDir) check() (status.Status, string, status.Detail) {
	f, err := os.Open(x.fpath)
	if err != nil {
		return status.StatusFailed, `Could not open "` + x.fpath + "`", status.Error{Err: err}
	}
	defer f.Close()

	currentConfig := map[string]string{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		matches := xdgUserDirRe.FindStringSubmatch(line)
		if len(matches) != 3 {
			continue
		}
		currentConfig[matches[1]] = strings.Replace(strings.TrimRight(matches[2], "/"), "$HOME", x.user.Home, 1)
	}

	if err := scanner.Err(); err != nil {
		x.msg(status.StatusFailed, `Could not read "`+x.fpath+`"`, status.Error{Err: err})
	}

	result := status.Table{Header: []string{"Directory", "Current", "Required"}}

	var todo bool
	for _, name := range xdgOrderedVars {
		dirStatus := status.StatusApplied
		if currentConfig[name] != x.required[name] {
			dirStatus = status.StatusTodo
			todo = true
		}

		result.AppendRow(
			status.TableCell{Status: dirStatus, Content: strings.ToTitle(name)},
			status.TableCell{Status: dirStatus, Content: currentConfig[name]},
			status.TableCell{Status: status.StatusNone, Content: x.required[name]},
		)
	}

	if todo {
		return status.StatusTodo, "Need to change XDG dirs for " + x.user.Name, &result
	}
	return status.StatusApplied, "XDG dirs for " + x.user.Name + " are as expected", &result
}

func (x *xdgUserDir) apply() bool {
	st, msg, det := x.check()
	if st != status.StatusTodo {
		x.msg(st, msg, det)
		return true
	}

	f, err := os.Create(x.fpath)
	if err != nil {
		x.msg(status.StatusFailed, `Could not open "`+x.fpath+`"`, status.Error{Err: err})
		return false
	}
	defer f.Close()

	for _, name := range xdgOrderedVars {
		path := x.required[name]
		if path == "" {
			path = "$HOME/"
		}
		if _, err := f.WriteString(
			"XDG_" + name + `_DIR="` + x.required[name] + "\"\n",
		); err != nil {
			x.msg(status.StatusFailed, `Could not write to "`+x.fpath+`"`, status.Error{Err: err})
			return false
		}
	}
	x.msg(status.StatusApplied, `Applied XDG user paths for `+x.user.Name, nil)
	return true
}
