package external

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/willoma/keepakonf/internal/log"
)

const (
	dpkgStatusDir         = "/var/lib/dpkg"
	dpkgStatusPath        = dpkgStatusDir + "/status"
	dpkgWatcherDedupDelay = 100 * time.Millisecond
)

type dpkgWatcher struct {
	receivers   map[chan<- map[string]DpkgPackage]struct{}
	receiversMu sync.Mutex

	packages   map[string]DpkgPackage
	packagesMu sync.Mutex
}

var (
	dpkgWatcherRunner     *dpkgWatcher
	dpkgWatcherRunnerOnce sync.Once
)

func (d *dpkgWatcher) run() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error(err, "Could not create file watcher for dpkg status")
		return
	}

	var (
		dedupTimer *time.Timer
		dedupMutex sync.Mutex
	)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Create) && event.Name == dpkgStatusPath {
					dedupMutex.Lock()
					if dedupTimer == nil {
						dedupTimer = time.AfterFunc(dpkgWatcherDedupDelay, func() {
							d.scan()

							d.packagesMu.Lock()
							defer d.packagesMu.Unlock()

							d.receiversMu.Lock()
							defer d.receiversMu.Unlock()

							for c := range d.receivers {
								c <- d.packages
							}
						})
					} else {
						dedupTimer.Reset(dpkgWatcherDedupDelay)
					}
					dedupMutex.Unlock()
				}

			case err, ok := <-watcher.Errors:
				log.Error(err, "Could not monitor dpkg status")
				if !ok {
					return
				}
			}
		}
	}()

	if err := watcher.Add(dpkgStatusDir); err != nil {
		log.Errorf(err, "Could not add directory %q to watcher", dpkgStatusDir)
	}

	d.scan()
}

func (d *dpkgWatcher) scan() {
	f, err := os.Open(dpkgStatusPath)
	if err != nil {
		log.Error(err, "Could not open dpkg status")
		return
	}
	defer f.Close()

	packages := map[string]DpkgPackage{}

	var (
		pkg       string
		version   string
		installed bool
	)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			// End of record, store it
			if pkg != "" {
				packages[pkg] = DpkgPackage{
					Name:      pkg,
					Version:   version,
					Installed: installed,
				}
			}
			installed = false
			pkg = ""
			version = ""
			continue
		}

		info := strings.SplitN(line, ": ", 2)
		if len(info) != 2 {
			continue
		}

		switch info[0] {
		case "Package":
			pkg = info[1]
		case "Version":
			version = info[1]
		case "Status":
			if info[1] == "install ok installed" {
				installed = true
			}
		}
	}

	if pkg != "" {
		packages[pkg] = DpkgPackage{
			Version:   version,
			Installed: installed,
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error(err, "Could not read dpkg status")
		return
	}

	d.packagesMu.Lock()
	d.packages = packages
	d.packagesMu.Unlock()
}

func (d *dpkgWatcher) listen() (target <-chan map[string]DpkgPackage, remove func()) {
	targetChan := make(chan map[string]DpkgPackage, 2)

	d.packagesMu.Lock()
	targetChan <- d.packages
	d.packagesMu.Unlock()

	d.receiversMu.Lock()
	d.receivers[targetChan] = struct{}{}
	d.receiversMu.Unlock()

	return targetChan, func() {
		d.receiversMu.Lock()
		delete(d.receivers, targetChan)
		d.receiversMu.Unlock()
		close(targetChan)
	}
}

func (d *dpkgWatcher) listPackages() map[string]DpkgPackage {
	d.packagesMu.Lock()
	defer d.packagesMu.Unlock()
	return d.packages
}

func initDpkgWatcher() {
	dpkgWatcherRunnerOnce.Do(func() {
		dpkgWatcherRunner = &dpkgWatcher{
			receivers: map[chan<- map[string]DpkgPackage]struct{}{},
			packages:  map[string]DpkgPackage{},
		}
		dpkgWatcherRunner.run()
	})
}

// DpkgListen returns the list of known packages whenever it changes.
func DpkgListen() (target <-chan map[string]DpkgPackage, remove func()) {
	initDpkgWatcher()

	return dpkgWatcherRunner.listen()
}

// DpkgPackages returns the list of known packages once.
func DpkgPackages() map[string]DpkgPackage {
	initDpkgWatcher()

	return dpkgWatcherRunner.listPackages()
}
