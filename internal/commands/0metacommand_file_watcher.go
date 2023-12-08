package commands

import (
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
)

type fileWatcherCommand interface {
	getPath() string
	newStatus(external.FileStatus)
	apply() bool
}

type fileWatcherCommandInit func(params map[string]any, msg status.SendStatus) fileWatcherCommand

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
		func(params map[string]any, msg status.SendStatus) Command {
			return &fileWatcher{cmd: cmdInit(params, msg)}
		},
	)
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
