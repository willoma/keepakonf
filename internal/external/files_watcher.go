package external

import (
	"errors"
	"fmt"
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
	FileStatusDirectory
	FileStatusNotFound
	fileStatusParentRemoved

	filesWatcherDedupDelay = 200 * time.Millisecond
)

var (
	filesWatcher = newFilesWatcherFront()

	errDirNotExist              = errors.New("directory does not exist")
	errFileAlreadyWatched       = errors.New("file already watched")
	errFileNotWatched           = errors.New("file not watched")
	errFileWatcherCouldNotWatch = errors.New("file watcher could not watch")
	errFileWatcherWrongTarget   = errors.New("wrong file watcher target")
	errNotADirectory            = errors.New("not a directory")
)

func getFileStatus(path string) FileStatus {
	finfo, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return FileStatusNotFound
		} else {
			log.Errorf(err, "Could not check info for watched file %q", path)
			return FileStatusUnknown
		}
	}
	if finfo.IsDir() {
		return FileStatusDirectory
	}
	return FileStatusFile
}

type filesWatcherBack struct {
	watcher      *fsnotify.Watcher
	mu           sync.Mutex
	targets      map[string]chan FileStatus
	targetsByDir map[string]map[chan FileStatus]struct{}
}

func newFilesWatcherBack() *filesWatcherBack {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err, "Could not create file watcher")
	}

	return &filesWatcherBack{
		watcher:      watcher,
		targets:      map[string]chan FileStatus{},
		targetsByDir: map[string]map[chan FileStatus]struct{}{},
	}
}

func (b *filesWatcherBack) watch() {
	for {
		select {
		case event, ok := <-b.watcher.Events:
			if !ok {
				return
			}

			b.forwardStatus(event)
		case err, ok := <-b.watcher.Errors:
			log.Error(err, "Could not monitor files")
			if !ok {
				return
			}
		}
	}
}

func (b *filesWatcherBack) forwardStatus(event fsnotify.Event) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if event.Has(fsnotify.Remove) {
		if targetDir, ok := b.targetsByDir[event.Name]; ok {
			for t := range targetDir {
				t <- fileStatusParentRemoved
			}
		}
	}

	target, ok := b.targets[event.Name]
	if !ok {
		return
	}

	switch {
	case event.Has(fsnotify.Create), event.Has(fsnotify.Write):
		target <- getFileStatus(event.Name)
	case event.Has(fsnotify.Remove), event.Has(fsnotify.Rename):
		target <- FileStatusNotFound
	}
}

func (b *filesWatcherBack) addTarget(path string) (chan FileStatus, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.targets[path]; ok {
		return nil, errFileAlreadyWatched
	}

	target := make(chan FileStatus)

	dirPath := filepath.Dir(path)
	if len(b.targetsByDir[dirPath]) == 0 {
		if err := b.watcher.Add(dirPath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil, errDirNotExist
			}
			return nil, fmt.Errorf("could not add directory %q to files watcher: %w", dirPath, err)
		}
		b.targetsByDir[dirPath] = map[chan FileStatus]struct{}{target: {}}
	} else {
		b.targetsByDir[dirPath][target] = struct{}{}
	}

	b.targets[path] = target

	firstValue := getFileStatus(path)
	go func() {
		target <- firstValue
	}()

	return target, nil
}

func (b *filesWatcherBack) removeTarget(path string, target chan FileStatus) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	currentTarget, ok := b.targets[path]
	if !ok {
		return errFileNotWatched
	}
	if currentTarget != target {
		return errFileWatcherWrongTarget
	}

	dirPath := filepath.Dir(path)

	delete(b.targets, path)

	delete(b.targetsByDir[dirPath], target)
	if len(b.targetsByDir[dirPath]) == 0 {
		if err := b.watcher.Remove(dirPath); err != nil {
			return fmt.Errorf("could not remove directory %q from files watcher: %w", dirPath, err)
		}
		delete(b.targetsByDir, dirPath)
	}

	return nil
}

type fileWatcherMiddle struct {
	path  string
	front *filesWatcherFront

	mu  sync.Mutex
	in  chan FileStatus
	out []chan FileStatus
	val FileStatus
}

func newFilesWatcherMiddle(path string, front *filesWatcherFront) *fileWatcherMiddle {
	m := &fileWatcherMiddle{
		path:  path,
		front: front,
	}
	go m.run()
	return m
}

func (m *fileWatcherMiddle) run() {
	for {
		err := m.watchFile()
		switch {
		case errors.Is(err, errDirNotExist):
			m.mu.Lock()
			for _, t := range m.out {
				t <- FileStatusNotFound
			}
			m.mu.Unlock()
			m.waitForParent()
		case errors.Is(err, errFileWatcherCouldNotWatch):
			// Try again later...
			time.Sleep(10 * time.Second)
		default:
			return
		}
	}
}

func (m *fileWatcherMiddle) watchFile() error {
	target, err := m.front.back.addTarget(m.path)
	if err != nil {
		if errors.Is(err, errDirNotExist) {
			return errDirNotExist
		}
		log.Errorf(err, "Could not watch for %q", m.path)
		return errFileWatcherCouldNotWatch
	}
	m.in = target

	timer := time.AfterFunc(filesWatcherDedupDelay, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, tgt := range m.out {
			tgt <- m.val
		}
	})

	for fstatus := range m.in {
		if fstatus == fileStatusParentRemoved {
			return errDirNotExist
		}
		if fstatus == FileStatusUnknown {
			continue
		}
		m.mu.Lock()
		m.val = fstatus
		timer.Reset(filesWatcherDedupDelay)
		m.mu.Unlock()
	}
	return nil
}

func (m *fileWatcherMiddle) waitForParent() {
	dir := filepath.Dir(m.path)
	tgt := m.front.subscribe(dir)
	for fstatus := range tgt {
		switch fstatus {
		case FileStatusDirectory:
			m.front.unsubscribe(dir, tgt)
			return
		case FileStatusFile:
			log.Errorf(errNotADirectory, "Could not wait on parent %q creation", dir)
		}
	}
}

func (m *fileWatcherMiddle) addTarget() chan FileStatus {
	m.mu.Lock()
	defer m.mu.Unlock()

	target := make(chan FileStatus, 2)
	target <- m.val
	m.out = append(m.out, target)

	return target
}

// removeTarget returns true if there is no more targets in this watcher
func (m *fileWatcherMiddle) removeTarget(target chan FileStatus) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, tgt := range m.out {
		if tgt == target {
			m.out[i] = m.out[len(m.out)-1]
			m.out = m.out[:len(m.out)-1]
		}
	}

	if len(m.out) == 0 {
		m.front.back.removeTarget(m.path, m.in)
		return true
	}

	return false
}

type filesWatcherFront struct {
	mu      sync.Mutex
	middles map[string]*fileWatcherMiddle
	back    *filesWatcherBack
}

func newFilesWatcherFront() *filesWatcherFront {
	back := newFilesWatcherBack()
	go back.watch()
	return &filesWatcherFront{
		middles: map[string]*fileWatcherMiddle{},
		back:    back,
	}
}

func (f *filesWatcherFront) subscribe(path string) chan FileStatus {
	f.mu.Lock()
	defer f.mu.Unlock()

	middle, ok := f.middles[path]
	if !ok {
		middle = newFilesWatcherMiddle(path, f)
		f.middles[path] = middle
	}
	return middle.addTarget()
}

func (f *filesWatcherFront) unsubscribe(path string, target chan FileStatus) {
	f.mu.Lock()
	defer f.mu.Unlock()

	middle, ok := f.middles[path]
	if !ok {
		return
	}

	if middle.removeTarget(target) {
		delete(f.middles, path)
	}
}

// WatchFile allows watching for files creation, change or removal, and
// differentiates files and directories.
func WatchFile(path string) (target <-chan FileStatus, remove func()) {
	targetChan := filesWatcher.subscribe(path)

	return targetChan, func() {
		filesWatcher.unsubscribe(path, targetChan)
	}
}
