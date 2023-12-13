package commands

import (
	"strings"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = register(
	"apt remove",
	"packages",
	"Remove packages using apt",
	[]ParamDesc{
		{"packages", "Packages to remove", ParamTypeStringArray},
		{"purge", "Purge the packages", ParamTypeBool},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) Command {
		var cmd string
		if params["purge"].(bool) {
			cmd = "purge"
		} else {
			cmd = "remove"
		}
		return &aptRemove{
			msg:      msg,
			vars:     vars,
			packages: params["packages"].([]string),
			cmd:      cmd,
		}
	},
)

type aptRemove struct {
	msg  status.SendStatus
	vars variables.Variables

	packages []string
	cmd      string

	needToRemove []string

	applying atomic.Bool
	close    func()
}

func (a *aptRemove) UpdateVariables(vars variables.Variables) {
	if a.vars.Update(vars) {
		a.update(external.DpkgPackages())
	}
}

func (a *aptRemove) Watch() {
	signals, close := external.DpkgListen()
	a.close = close

	a.update(external.DpkgPackages())

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
	pkgs := a.vars.ReplaceSlice(a.packages)

	for _, pkg := range pkgs {
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
	} else if len(pkgs) == 1 {
		info = "Package " + pkgs[0] + " removed"
	} else {
		info = "Packages " + strings.Join(pkgs, ", ") + " removed"
	}

	a.needToRemove = needToRemove
	if !a.applying.Load() {
		a.msg(msgStatus, info, &table, nil)
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
		func(s status.Status, info string, detail status.Detail) {
			if info == "" {
				switch s {
				case status.StatusRunning:
					info = "Removing " + needToRemoveMsg
				case status.StatusApplied:
					info = "Successfully removed " + needToRemoveMsg
				case status.StatusFailed:
					info = "Failed removing " + needToRemoveMsg
				}

			}
			a.msg(s, info, detail, nil)
		},
		a.cmd, a.needToRemove,
	)
}
