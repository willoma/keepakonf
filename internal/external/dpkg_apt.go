package external

import (
	"sync"

	"github.com/willoma/keepakonf/internal/status"
)

var dpkgMu sync.Mutex

func AptGet(
	receiver func(status.Status, string, status.Detail),
	cmd string,
	packages []string,
) bool {
	if !dpkgMu.TryLock() {
		// Try locking only to send the message if the mutex cannot be locked
		// right now. It allows not sending the message if dpkg is already
		// available...
		receiver(status.StatusRunning, "Waiting for dpkg to be available", nil)
		dpkgMu.Lock()
	}
	defer dpkgMu.Unlock()

	return execToMessage(
		receiver,
		[]string{"DEBIAN_FRONTEND=noninteractive"},
		"apt-get",
		append([]string{"--yes", "--quiet", cmd}, packages...)...,
	)
}
