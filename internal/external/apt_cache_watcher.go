package external

import (
	"bufio"
	"io"
	"log/slog"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
)

const (
	// This directory is modified when apt updates the list of packages
	aptUpdateFileToWatch = "/var/lib/apt/lists/partial"

	// Interval between two forced updates
	aptUpdateInterval = 2 * time.Hour

	// Some lines may contain more than the default of 64k characters.
	// (ex. Package librust-winapi-dev, entry "Provides": more than 70k)
	aptUpdateScannerBufferLength = 128 * 1024
)

type aptCacheWatcher struct {
	receivers     map[chan<- map[string]DpkgPackage]struct{}
	receiversList map[chan<- []DpkgPackage]struct{}
	receiversMu   sync.Mutex

	packages     map[string]DpkgPackage
	packagesList []DpkgPackage
	packagesMu   sync.Mutex
}

var (
	aptCacheWatcherRunner     *aptCacheWatcher
	aptCacheWatcherRunnerOnce sync.Once
)

func (a *aptCacheWatcher) run() {
	fstatusChan, _ := WatchFile(aptUpdateFileToWatch)
	ticker := time.NewTicker(aptUpdateInterval)

	knownPackages, _ := DpkgListen()

	go func() {
		for {
			select {
			case fstatus := <-fstatusChan:
				if fstatus == FileStatusDirectory {
					// Seems like "apt update" has been executed...

					a.updatePackagesList(DpkgPackages())
					ticker.Reset(aptUpdateInterval)
				}
			case currentPackages := <-knownPackages:
				a.updatePackagesList(currentPackages)

			case <-ticker.C:
				// Force an update on a regular basis...

				// Don't lock the mutex just to check the number of receivers,
				// it's a little bit risky, but race conditions will probably
				// never happen.
				if len(a.receivers) == 0 {
					// Ignore if there is no listener...
					continue
				}

				if AptGet(
					func(s status.Status, info string, detail status.Detail) {
						if s == status.StatusFailed {
							log.Info("Could not download apt packages list", "error", status.StatusFailed, "", "", "", status.DetailJSON(detail))
						}
					},
					"update",
				) {
					a.updatePackagesList(DpkgPackages())
				}
			}
		}
	}()
}

func (a *aptCacheWatcher) updatePackagesList(knownPackages map[string]DpkgPackage) {
	cmd := exec.Command("apt-cache", "dumpavail")
	reader, writer := io.Pipe()
	cmd.Stdout = writer
	if err := cmd.Start(); err != nil {
		log.Error(err, "Could not get information from apt cache")
		return
	}

	packages := map[string]DpkgPackage{}

	go func() {
		var (
			pkg     string
			version string
		)

		buf := []byte{}
		scanner := bufio.NewScanner(reader)
		scanner.Buffer(buf, aptUpdateScannerBufferLength)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if line == "" {
				// End of record, compare it and store it
				if pkg != "" {
					pkgObj := DpkgPackage{
						Name:             pkg,
						AvailableVersion: version,
					}
					known, ok := knownPackages[pkg]
					if ok && known.Installed {
						pkgObj.Version = known.Version
						pkgObj.Installed = true
					}
					if _, ok := packages[pkg]; !ok || packages[pkg].AvailableVersion < pkgObj.AvailableVersion {
						packages[pkg] = pkgObj
					}
				}

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
			}
		}

		if err := scanner.Err(); err != nil {
			slog.Error("err", "error", err)
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Error(err, "Could not read information from apt cache")
		return
	}

	packagesList := make([]DpkgPackage, 0, len(packages))
	for _, pkg := range packages {
		packagesList = append(packagesList, pkg)
	}

	sort.Slice(packagesList, func(i, j int) bool {
		return packagesList[i].Name < packagesList[j].Name
	})

	a.packagesMu.Lock()
	a.packages = packages
	a.packagesList = packagesList
	a.packagesMu.Unlock()

	a.receiversMu.Lock()
	for c := range a.receivers {
		c <- a.packages
	}
	for c := range a.receiversList {
		c <- a.packagesList
	}
	a.receiversMu.Unlock()
}

func (a *aptCacheWatcher) listen() (target <-chan map[string]DpkgPackage, remove func()) {
	targetChan := make(chan map[string]DpkgPackage)
	a.receiversMu.Lock()
	a.receivers[targetChan] = struct{}{}
	a.receiversMu.Unlock()

	return targetChan, func() {
		a.receiversMu.Lock()
		delete(a.receivers, targetChan)
		a.receiversMu.Unlock()
		close(targetChan)
	}
}

func (a *aptCacheWatcher) listenList() (target <-chan []DpkgPackage, remove func()) {
	targetChan := make(chan []DpkgPackage)
	a.receiversMu.Lock()
	a.receiversList[targetChan] = struct{}{}
	a.receiversMu.Unlock()

	return targetChan, func() {
		a.receiversMu.Lock()
		delete(a.receiversList, targetChan)
		a.receiversMu.Unlock()
		close(targetChan)
	}
}

func initAptCacheWatcher() {
	aptCacheWatcherRunnerOnce.Do(func() {
		aptCacheWatcherRunner = &aptCacheWatcher{
			receivers:     map[chan<- map[string]DpkgPackage]struct{}{},
			receiversList: map[chan<- []DpkgPackage]struct{}{},
			packages:      map[string]DpkgPackage{},
		}
		aptCacheWatcherRunner.run()
	})
}

// AptCacheListen returns the list of known packages and their potential update whenever it changes.
func AptCacheListen() (target <-chan map[string]DpkgPackage, remove func()) {
	initAptCacheWatcher()

	return aptCacheWatcherRunner.listen()
}

// AptCacheListen returns the list of known packages and their potential update whenever it changes.
func AptCacheListenList() (target <-chan []DpkgPackage, remove func()) {
	initAptCacheWatcher()

	return aptCacheWatcherRunner.listenList()
}
