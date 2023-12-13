package commands

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
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
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCommand {
		return &xdgUserDir{
			msg:  msg,
			vars: vars,
			user: params["user"].(string),
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
	msg  status.SendStatus
	vars variables.Variables

	user string

	required map[string]string
}

func (x *xdgUserDir) updateVariables(vars variables.Variables) (changed bool) {
	return x.vars.Update(vars)
}

func (x *xdgUserDir) getPath() string {
	userData, err := external.GetUser(x.vars.Replace(x.user))
	if err != nil {
		log.Errorf(err, "Could not get user information for %q", x.vars.Replace(x.user))
		return filepath.Join(os.TempDir(), xdgUserDirPath)
	}

	return filepath.Join(userData.Home, xdgUserDirPath)
}

func (x *xdgUserDir) newStatus(fstatus external.FileStatus) {
	switch fstatus {
	case external.FileStatusFile:
		x.msg(x.check())
	case external.FileStatusDirectory:
		x.msg(status.StatusFailed, `"`+x.getPath()+`" is a directory`, nil, nil)
	case external.FileStatusNotFound:
		x.msg(status.StatusFailed, `File "`+x.getPath()+`" not found`, nil, nil)
	}
}

func (x *xdgUserDir) check() (status.Status, string, status.Detail, variables.Variables) {
	userData, err := external.GetUser(x.vars.Replace(x.user))
	if err != nil {
		log.Errorf(err, "Could not get user information for %q", x.vars.Replace(x.user))
	}

	f, err := os.Open(x.getPath())
	if err != nil {
		return status.StatusFailed, `Could not open "` + x.getPath() + "`", status.Error(err.Error()), nil
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
		currentConfig[matches[1]] = strings.Replace(strings.TrimRight(matches[2], "/"), "$HOME", userData.Home, 1)
	}

	if err := scanner.Err(); err != nil {
		x.msg(status.StatusFailed, `Could not read "`+x.getPath()+`"`, status.Error(err.Error()), nil)
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
		return status.StatusTodo, "Need to change XDG dirs for " + userData.Name, &result, nil
	}
	return status.StatusApplied, "XDG dirs for " + userData.Name + " are as expected", &result, nil
}

func (x *xdgUserDir) apply() bool {
	st, msg, det, vars := x.check()
	if st != status.StatusTodo {
		x.msg(st, msg, det, vars)
		return true
	}

	f, err := os.Create(x.getPath())
	if err != nil {
		x.msg(status.StatusFailed, `Could not open "`+x.getPath()+`"`, status.Error(err.Error()), nil)
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
			x.msg(status.StatusFailed, `Could not write to "`+x.getPath()+`"`, status.Error(err.Error()), nil)
			return false
		}
	}

	// TODO change ownership

	x.msg(status.StatusApplied, `Applied XDG user paths for `+x.vars.Replace(x.user), nil, nil)
	return true
}
