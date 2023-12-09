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
	FileStatusDirectory
	FileStatusNotFound

	filesWatcherDedupDelay = 200 * time.Millisecond
)

type filesWatcherStruct struct {
	once     sync.Once
	watcher  *fsnotify.Watcher
	mu       sync.Mutex
	dedups   map[string]*fileStatusDedup
	dirCount map[string]int
}

type fileStatusDedup struct {
	in  chan FileStatus
	mu  sync.Mutex
	out []chan FileStatus
	val FileStatus
}

var filesWatcher = filesWatcherStruct{
	dedups:   map[string]*fileStatusDedup{},
	dirCount: map[string]int{},
}

func (d *fileStatusDedup) listen() {
	go func() {
		timer := time.AfterFunc(filesWatcherDedupDelay, func() {
			d.mu.Lock()
			defer d.mu.Unlock()
			for _, r := range d.out {
				r <- d.val
			}
		})

		for fstatus := range d.in {
			if fstatus == FileStatusUnknown {
				continue
			}
			d.mu.Lock()
			d.val = fstatus
			timer.Reset(filesWatcherDedupDelay)
			d.mu.Unlock()
		}
	}()
}
func (d *fileStatusDedup) addReceiver() chan FileStatus {
	d.mu.Lock()
	defer d.mu.Unlock()
	receiver := make(chan FileStatus)
	d.out = append(d.out, receiver)
	firstValue := d.val
	go func() {
		receiver <- firstValue
	}()
	return receiver
}

func (d *fileStatusDedup) removeReceiver(r chan FileStatus) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i, rec := range d.out {
		if rec == r {
			d.out[i] = d.out[len(d.out)-1]
			d.out = d.out[:len(d.out)-1]
		}
	}
	return len(d.out) == 0
}

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

func (f *filesWatcherStruct) run() {
	f.once.Do(func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Error(err, "Could not create file watcher")
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

					f.mu.Lock()
					target, ok := f.dedups[event.Name]
					f.mu.Unlock()
					if !ok {
						continue
					}

					switch {
					case event.Has(fsnotify.Create), event.Has(fsnotify.Write):
						target.in <- getFileStatus(event.Name)
					case event.Has(fsnotify.Remove), event.Has(fsnotify.Rename):
						target.in <- FileStatusNotFound
					default:
						// Ignore chmod events
					}
				case err, ok := <-watcher.Errors:
					log.Error(err, "Could not monitor files")
					if !ok {
						return
					}
				}
			}
		}()
	})
}

func (f *filesWatcherStruct) subscribe(path string) chan FileStatus {
	f.mu.Lock()
	defer f.mu.Unlock()
	dedup, ok := f.dedups[path]
	if !ok {
		dedup = &fileStatusDedup{
			in:  make(chan FileStatus),
			val: getFileStatus(path),
		}

		dir := filepath.Dir(path)
		if f.dirCount[dir] == 0 {
			if err := f.watcher.Add(dir); err != nil {
				// TODO If directory does not exist, listen for its parent, etc
				log.Errorf(err, "Could not add directory %q to files watcher", dir)
			}
		}
		dedup.listen()
		f.dirCount[dir]++
		f.dedups[path] = dedup
	}

	return dedup.addReceiver()
}

func (f *filesWatcherStruct) unsubscribe(path string, receiver chan FileStatus) {
	f.mu.Lock()
	defer f.mu.Unlock()
	close(receiver)
	dedup, ok := f.dedups[path]
	if !ok {
		return
	}
	if dedup.removeReceiver(receiver) {
		close(dedup.in)
		delete(f.dedups, path)
	}
	dir := filepath.Dir(path)
	f.dirCount[dir]--
	if f.dirCount[dir] == 0 {
		if err := f.watcher.Remove(dir); err != nil {
			log.Errorf(err, "Could not remove directory %q from files watcher", dir)
		}
	}
}

// WatchFile allows watching for files creation, change or removal, and
// differentiates files and directories.
func WatchFile(path string) (target <-chan FileStatus, remove func()) {
	filesWatcher.run()

	targetChan := filesWatcher.subscribe(path)

	return targetChan, func() {
		filesWatcher.unsubscribe(path, targetChan)
	}
}
