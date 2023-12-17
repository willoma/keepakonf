package external

import "sync"

type DpkgPackage struct {
	Name             string
	Installed        bool
	Version          string
	AvailableVersion string
}

var dpkgMu sync.Mutex
