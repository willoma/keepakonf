package commands

import (
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

type fileWatcherCommand interface {
	updateVariables(vars variables.Variables) (changed bool)
	getPath() string
	newStatus(external.FileStatus)
	apply() bool
}

type fileWatcherCmd struct {
	msg  status.SendStatus
	vars variables.Variables

	path string
}

func fileWatcherCmdInit(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCmd {
	return fileWatcherCmd{
		msg:  msg,
		vars: vars,
		path: params["path"].(string),
	}
}

func (f *fileWatcherCmd) updateVariables(vars variables.Variables) (changed bool) {
	return f.vars.Update(vars)
}

func (f *fileWatcherCmd) getPath() string {
	return f.vars.Replace(f.path)
}

type fileWatcherCommandInit func(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCommand

type fileWatcher struct {
	cmd fileWatcherCommand

	applying atomic.Bool
	close    func()
}

func registerFileWatcher(
	name, icon, description string,
	parameters ParamsDesc,
	cmdInit fileWatcherCommandInit,
) struct{} {
	return register(
		name, icon, description, parameters,
		func(params map[string]any, vars variables.Variables, msg status.SendStatus) Command {
			return &fileWatcher{cmd: cmdInit(params, vars, msg)}
		},
	)
}

func (f *fileWatcher) UpdateVariables(vars variables.Variables) {
	if f.cmd.updateVariables(vars) {
		f.Stop()
		f.Watch()
	}
}

func (f *fileWatcher) Watch() {
	signals, close := external.WatchFile(f.cmd.getPath())
	f.close = close

	go func() {
		for fstatus := range signals {
			if f.applying.Load() {
				// No update if it is currently applying
				continue
			}

			f.cmd.newStatus(fstatus)
		}
	}()
}
func (f *fileWatcher) Stop() {
	if f.close != nil {
		f.close()
	}
}

func (f *fileWatcher) Apply() bool {
	f.applying.Store(true)
	defer f.applying.Store(false)
	return f.cmd.apply()
}
