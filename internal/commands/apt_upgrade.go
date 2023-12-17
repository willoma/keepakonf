package commands

import (
	"sync/atomic"
	"time"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = register(
	"apt upgrade",
	"packages",
	"Upgrade packages from APT repositories",
	ParamsDesc{},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) Command {
		return &aptUpgrade{
			msg:  msg,
			vars: vars,
		}
	},
)

type aptUpgrade struct {
	msg  status.SendStatus
	vars variables.Variables

	applying  atomic.Bool
	close     func()
	closeChan chan struct{}
	ticker    *time.Ticker
}

func (a *aptUpgrade) UpdateVariables(vars variables.Variables) {
	a.vars.Update(vars)
}

func (a *aptUpgrade) Watch() {
	aptcache, close := external.AptCacheListenList()
	a.close = close

	go func() {
		for packages := range aptcache {
			if a.applying.Load() {
				// No update if it is currently applying
				continue
			}

			table := status.Table{
				Header: []string{"Package", "Installed version", "Available version"},
			}

			var needToUpgrade bool

			for _, pkg := range packages {
				if !pkg.Installed || pkg.Version == pkg.AvailableVersion {
					continue
				}

				table.AppendRow(
					status.TableCell{Status: status.StatusNone, Content: pkg.Name},
					status.TableCell{Status: status.StatusNone, Content: pkg.Version},
					status.TableCell{Status: status.StatusNone, Content: pkg.AvailableVersion},
				)
				needToUpgrade = true
			}

			if needToUpgrade {
				a.msg(status.StatusTodo, "Need to upgrade packages", &table, nil)
			} else {
				a.msg(status.StatusApplied, "All packages up-to-date", nil, nil)
			}
		}
	}()
}

func (a *aptUpgrade) Stop() {
	if a.close != nil {
		a.close()
	}
	if a.ticker != nil {
		a.ticker.Stop()
	}
	a.closeChan <- struct{}{}
}

func (a *aptUpgrade) Apply() bool {
	a.applying.Store(true)
	defer a.applying.Store(false)

	return external.AptGet(
		func(s status.Status, info string, detail status.Detail) {
			if info == "" {
				switch s {
				case status.StatusRunning:
					info = "Upgrading packages"
				case status.StatusApplied:
					info = "Successfully upgraded packages"
				case status.StatusFailed:
					info = "Failed upgrading packages"
				}

			}
			a.msg(s, info, detail, nil)
		},
		"upgrade",
	)
}
