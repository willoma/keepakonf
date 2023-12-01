package external

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/willoma/keepakonf/internal/log"
)

type FileStatus int

const (
	FileStatusUnknown FileStatus = iota
	FileStatusFile
	FileStatusDirectory
	FileStatusNotFound
)

type filesWatcher struct {
	logger            *log.Logger
	filePathReceivers map[string][]chan FileStatus
	dirCount          map[string]int
	receiversMu       sync.Mutex

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

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				f.newEvent(event)
			case err, ok := <-watcher.Errors:
				f.logger.Error(err, "Could not monitor files")
				if !ok {
					return
				}
			}
		}
	}()
}

func (f *filesWatcher) newEvent(event fsnotify.Event) {
	f.receiversMu.Lock()
	defer f.receiversMu.Unlock()

	receivers, ok := f.filePathReceivers[event.Name]
	if !ok {
		return
	}

	switch {
	case event.Has(fsnotify.Create):
		finfo, err := os.Stat(event.Name)
		if err != nil {
			f.logger.Errorf(err, "Could not check info for watched file %q", event.Name)
			return
		}
		if finfo.IsDir() {
			for _, r := range receivers {
				r <- FileStatusDirectory
			}
		} else {
			for _, r := range receivers {
				r <- FileStatusFile
			}
		}
	case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
		for _, r := range receivers {
			r <- FileStatusNotFound
		}
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

func WatchFile(logger *log.Logger, filePath string) (target <-chan FileStatus, remove func()) {
	initFilesWatcher(logger)

	return filesWatcherRunner.listen(filePath)
}
