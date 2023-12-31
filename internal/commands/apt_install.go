package commands

import (
	"strings"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

var _ = register(
	"apt install",
	"packages",
	"Install packages using apt",
	ParamsDesc{
		{"packages", "Packages to install", ParamTypeStringArray},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) Command {
		return &aptInstall{
			msg:      msg,
			vars:     vars,
			packages: params["packages"].([]string),
		}
	},
)

type aptInstall struct {
	msg  status.SendStatus
	vars variables.Variables

	packages []string

	needToInstall []string

	applying atomic.Bool
	close    func()
}

func (a *aptInstall) UpdateVariables(vars variables.Variables) {
	if a.vars.Update(vars) {
		a.update(external.AptCachePackages())
	}
}

func (a *aptInstall) Watch() {
	signals, close := external.AptCacheListen()
	a.close = close

	go func() {
		for knownPackages := range signals {
			a.update(knownPackages)
		}
	}()
}

func (a *aptInstall) update(knownPackages map[string]external.DpkgPackage) {
	needToInstall := []string{}
	unknown := []string{}

	msgStatus := status.StatusApplied
	table := status.Table{
		Header: []string{"Package", "Installed Version", "Available Version"},
	}

	pkgs := a.vars.ReplaceSlice(a.packages)

	for _, pkg := range pkgs {
		if pkg == "" {
			continue
		}
		info, ok := knownPackages[pkg]
		switch {
		case !ok:
			unknown = append(unknown, pkg)
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusFailed, Content: "None"},
				status.TableCell{Status: status.StatusNone, Content: "Unknown"},
			)
		case info.Installed:
			switch {
			case info.AvailableVersion == "":
				table.AppendRow(
					status.TableCell{Status: status.StatusNone, Content: pkg},
					status.TableCell{Status: status.StatusApplied, Content: info.Version},
					status.TableCell{Status: status.StatusNone, Content: "None"},
				)
			case info.Version == info.AvailableVersion:
				table.AppendRow(
					status.TableCell{Status: status.StatusNone, Content: pkg},
					status.TableCell{Status: status.StatusApplied, Content: info.Version},
					status.TableCell{Status: status.StatusNone, Content: info.AvailableVersion},
				)
			default:
				needToInstall = append(needToInstall, pkg)
				table.AppendRow(
					status.TableCell{Status: status.StatusNone, Content: pkg},
					status.TableCell{Status: status.StatusTodo, Content: info.Version},
					status.TableCell{Status: status.StatusNone, Content: info.AvailableVersion},
				)
			}
		default:
			needToInstall = append(needToInstall, pkg)
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusTodo, Content: "None"},
				status.TableCell{Status: status.StatusNone, Content: info.AvailableVersion},
			)
		}
	}

	var info string
	if len(unknown) > 0 {
		msgStatus = status.StatusFailed
		info = strings.Join(needToInstall, ", ") + "unknown"
	} else if len(needToInstall) > 0 {
		msgStatus = status.StatusTodo
		info = "Need to install " + strings.Join(needToInstall, ", ")
	} else if len(pkgs) == 1 {
		info = "Package " + pkgs[0] + " installed"
	} else {
		info = "Packages " + strings.Join(pkgs, ", ") + " installed"
	}

	a.needToInstall = needToInstall
	if !a.applying.Load() {
		a.msg(msgStatus, info, &table, nil)
	}
}

func (a *aptInstall) Stop() {
	if a.close != nil {
		a.close()
	}
}

func (a *aptInstall) Apply() bool {
	needToInstallMsg := strings.Join(a.needToInstall, ", ")

	a.applying.Store(true)
	defer a.applying.Store(false)

	return external.AptGet(
		func(s status.Status, info string, detail status.Detail) {
			if info == "" {
				switch s {
				case status.StatusRunning:
					info = "Installing " + needToInstallMsg
				case status.StatusApplied:
					info = "Successfully installed " + needToInstallMsg
				case status.StatusFailed:
					info = "Failed installing " + needToInstallMsg
				}

			}
			a.msg(s, info, detail, nil)
		},
		"install", a.needToInstall...,
	)
}
