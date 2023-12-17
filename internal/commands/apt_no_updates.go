package commands

import (
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

const aptNoUpdates = `# Placed by Quickonf
APT::Periodic::Update-Package-Lists "0";
APT::Periodic::Download-Upgradeable-Packages "0";
APT::Periodic::AutocleanInterval "0";
`

var _ = registerFileWatcher(
	"apt no updates",
	"ubuntu",
	"Disable periodic apt update, upgrade, autoclean",
	ParamsDesc{},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) fileWatcherCommand {
		return &fileContent{
			fileWatcherCmdInit(
				map[string]any{"path": "/etc/apt/apt.conf.d/10periodic"},
				vars, msg,
			),
			aptNoUpdates,
			"root",
		}
	},
)
