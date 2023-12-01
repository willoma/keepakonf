package commands

import (
	"strings"
	"sync/atomic"

	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/log"
	"github.com/willoma/keepakonf/internal/status"
)

var _ = register(
	"apt install",
	"packages",
	"Install packages using apt",
	ParamsDesc{
		{"packages", "Packages to install", ParamTypeStringArray},
	},
	func(params map[string]any, logger *log.Logger, msg status.SendStatus) Command {
		return &aptInstall{
			command:  command{logger, msg},
			packages: params["packages"].([]string),
		}
	},
)

type aptInstall struct {
	command

	packages []string

	needToInstall []string

	applying atomic.Bool
	close    func()
}

func (a *aptInstall) Watch() {
	signals, close := external.DpkgListen(a.logger)
	a.close = close

	a.update(external.DpkgPackages(a.logger))

	go func() {
		for knownPackages := range signals {
			a.update(knownPackages)
		}
	}()
}

func (a *aptInstall) update(knownPackages map[string]external.DpkgPackage) {
	needToInstall := []string{}

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
			needToInstall = append(needToInstall, pkg)
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusUnknown, Content: "not installed"},
				status.TableCell{Status: status.StatusNone, Content: ""},
			)
		case info.Installed:
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusApplied, Content: "installed"},
				status.TableCell{Status: status.StatusNone, Content: info.Version},
			)
		default:
			needToInstall = append(needToInstall, pkg)
			table.AppendRow(
				status.TableCell{Status: status.StatusNone, Content: pkg},
				status.TableCell{Status: status.StatusUnknown, Content: "not installed"},
				status.TableCell{Status: status.StatusNone, Content: info.Version},
			)
		}
	}

	var info string
	if len(needToInstall) > 0 {
		msgStatus = msgStatus.IfHigherPriority(status.StatusTodo)
		info = "Need to install " + strings.Join(needToInstall, ", ")
	} else if len(a.packages) == 1 {
		info = "Package " + a.packages[0] + " installed"
	} else {
		info = "Packages " + strings.Join(a.packages, ", ") + " installed"
	}

	a.needToInstall = needToInstall
	if !a.applying.Load() {
		a.msg(msgStatus, info, &table)
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
		"Installing "+needToInstallMsg,
		"Successfully installed "+needToInstallMsg,
		"Failed installing "+needToInstallMsg,
		a.msg,
		"install", a.needToInstall,
	)
}
