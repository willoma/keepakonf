package external

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/willoma/keepakonf/internal/log"
)

type FileStatus int

const (
	FileStatusUnknown FileStatus = iota
	FileStatusFile
	FileStatusFileChange
	FileStatusDirectory
	FileStatusNotFound

	filesWatcherDedupDelay = 200 * time.Millisecond
)

type filesWatcher struct {
	logger            *log.Logger
	filePathReceivers map[string][]chan FileStatus
	dirCount          map[string]int
	receiversMu       sync.Mutex

	newerEvents   map[string][]fsnotify.Op
	newerEventsMu sync.Mutex

	watcher *fsnotify.Watcher
}

var (
	filesWatcherRunner     *filesWatcher
	filesWatcherRunnerOnce sync.Once
)

func (f *filesWatcher) run() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		f.logger.Error(err, "Could not create file watcher for files")
		return
	}
	f.watcher = watcher

	var (
		dedupTimers map[string]*time.Timer
		dedupMutex  sync.Mutex
	)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				f.newerEventsMu.Lock()
				f.newerEvents[event.Name] = append(f.newerEvents[event.Name], event.Op)
				f.newerEventsMu.Unlock()

				dedupMutex.Lock()
				if dedupTimers[event.Name] == nil {
					dedupTimers[event.Name] = time.AfterFunc(filesWatcherDedupDelay, func() {
						f.forwardEvents(event.Name)
					})
				} else {
					dedupTimers[event.Name].Reset(filesWatcherDedupDelay)
				}
				dedupMutex.Unlock()
			case err, ok := <-watcher.Errors:
				f.logger.Error(err, "Could not monitor files")
				if !ok {
					return
				}
			}
		}
	}()
}

func (f *filesWatcher) forwardEvents(name string) {
	f.receiversMu.Lock()
	defer f.receiversMu.Unlock()

	f.newerEventsMu.Lock()
	defer f.newerEventsMu.Unlock()

	receivers, ok := f.filePathReceivers[name]
	if !ok {
		return
	}

	var previousEv fsnotify.Op = 0
	for _, ev := range f.newerEvents[name] {
		if ev == previousEv {
			continue
		}
		status := f.forwardEvent(name, ev)
		for _, r := range receivers {
			r <- status
		}
		previousEv = ev
	}
}

func (f *filesWatcher) forwardEvent(name string, op fsnotify.Op) FileStatus {
	switch {
	case op.Has(fsnotify.Create):
		finfo, err := os.Stat(name)
		if err != nil {
			f.logger.Errorf(err, "Could not check info for watched file %q", name)
			return FileStatusUnknown
		}
		if finfo.IsDir() {
			return FileStatusDirectory
		} else {
			return FileStatusFile
		}
	case op.Has(fsnotify.Write):
		return FileStatusFileChange
	case op.Has(fsnotify.Remove), op.Has(fsnotify.Rename):
		return FileStatusNotFound
	default:
		return FileStatusUnknown
	}
}

func (f *filesWatcher) startWatching(dirName string) {
	if err := f.watcher.Add(dirName); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			f.logger.Errorf(err, "Could not add directory %q to files watcher because its parent does not exist", dirName)
			return
		}
		f.logger.Errorf(err, "Could not add directory %q to files watcher", dirName)
	}
}

func (f *filesWatcher) stopWatching(dirName string) {
	if err := f.watcher.Remove(dirName); err != nil {
		f.logger.Errorf(err, "Could not remove directory %q from files watcher", dirName)
	}
}

func (f *filesWatcher) listen(path string) (target <-chan FileStatus, remove func()) {
	dirName := filepath.Dir(path)
	targetChan := make(chan FileStatus)

	f.receiversMu.Lock()

	f.filePathReceivers[path] = append(f.filePathReceivers[path], targetChan)
	if f.dirCount[dirName] == 0 {
		f.startWatching(dirName)
	}
	f.dirCount[dirName]++

	f.receiversMu.Unlock()

	var firstStatus FileStatus
	if finfo, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			firstStatus = FileStatusNotFound
		} else {
			f.logger.Errorf(err, "Could not check %q status before watching", path)
		}
	} else if finfo.IsDir() {
		firstStatus = FileStatusDirectory
	} else {
		firstStatus = FileStatusFile
	}

	go func() {
		targetChan <- firstStatus
	}()

	return targetChan, func() {
		f.receiversMu.Lock()

		newReceivers := make([]chan FileStatus, len(f.filePathReceivers[path]))
		for _, r := range f.filePathReceivers[path] {
			if r != targetChan {
				newReceivers = append(newReceivers, r)
			}
		}
		if len(newReceivers) == 0 {
			delete(f.filePathReceivers, path)
		} else {
			f.filePathReceivers[path] = newReceivers
		}

		f.dirCount[dirName]--
		if f.dirCount[dirName] == 0 {
			f.stopWatching(dirName)
		}

		f.receiversMu.Unlock()
		close(targetChan)
	}
}

func initFilesWatcher(logger *log.Logger) {
	filesWatcherRunnerOnce.Do(func() {
		filesWatcherRunner = &filesWatcher{
			filePathReceivers: map[string][]chan FileStatus{},
			dirCount:          map[string]int{},
			logger:            logger,
		}
		filesWatcherRunner.run()
	})
}

// WatchFile allows watching for files creation, change or removal, and
// differentiates files and directories on creation.
func WatchFile(logger *log.Logger, filePath string) (target <-chan FileStatus, remove func()) {
	initFilesWatcher(logger)

	return filesWatcherRunner.listen(filePath)
}
