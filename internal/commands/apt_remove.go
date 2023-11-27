package commands

import (
	"strings"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
)

var _ = register(
	"apt remove",
	"packages",
	"Remove packages using apt",
	[]ParamDesc{
		{"packages", "Packages to remove", ParamTypeStringArray},
		{"purge", "Purge the packages", ParamTypeBool},
	},
	func(params map[string]any, logger *log.Logger, msg status.SendStatus) Command {
		var command string
		if params["purge"].(bool) {
			command = "purge"
		} else {
			command = "remove"
		}
		return &aptRemove{
			packages: params["packages"].([]string),
			command:  command,
			logger:   logger,
			msg:      msg,
		}
	},
)

type aptRemove struct {
	packages []string
	command  string

	needToRemove []string

	logger *log.Logger
	msg    status.SendStatus

	applying atomic.Bool

	close func()
}

func (a *aptRemove) Watch() {
	signals, close := external.DpkgListen(a.logger)
	a.close = close

	a.update(external.DpkgPackages(a.logger))

	go func() {
		for knownPackages := range signals {
			a.update(knownPackages)
		}
	}()

}

func (a *aptRemove) update(knownPackages map[string]external.DpkgPackage) {
	needToRemove := []string{}

	msgStatus := status.StatusApplied
	table := status.Table{
		Header: []string{"Package", "Status", "Version"},
	}

	for _, pkg := range a.packages {
		if pkg == "" {
			continue
		}
		info, ok := knownPackages[pkg]
		switch {
		case !ok:
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusApplied, Content: "not installed"},
				status.TableCell{Status: status.StatusNone, Content: ""},
			)
		case info.Installed:
			needToRemove = append(needToRemove, pkg)
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusTodo, Content: "installed"},
				status.TableCell{Status: status.StatusNone, Content: info.Version},
			)
		default:
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusApplied, Content: "not installed"},
				status.TableCell{Status: status.StatusNone, Content: info.Version},
			)
		}
	}

	var info string
	if len(needToRemove) > 0 {
		msgStatus = msgStatus.IfHigherPriority(status.StatusTodo)
		info = "Need to remove " + strings.Join(needToRemove, ", ")
	} else if len(a.packages) == 1 {
		info = "Package " + a.packages[0] + " removed"
	} else {
		info = "Packages " + strings.Join(a.packages, ", ") + " removed"
	}

	a.needToRemove = needToRemove
	if !a.applying.Load() {
		a.msg(msgStatus, info, &table)
	}
}

func (a *aptRemove) Stop() {
	if a.close != nil {
		a.close()
	}
}

func (a *aptRemove) Apply() bool {
	needToRemoveMsg := strings.Join(a.needToRemove, ", ")

	a.applying.Store(true)
	defer a.applying.Store(false)

	return external.AptGet(
		"Removing "+needToRemoveMsg,
		"Successfully removed "+needToRemoveMsg,
		"Failed removing "+needToRemoveMsg,
		a.msg,
		a.command, a.needToRemove,
	)
}
