package commands

import (
	"github.com/willoma/keepakonf/internal/external"
	"github.com/willoma/keepakonf/internal/status"
	"github.com/willoma/keepakonf/internal/variables"
)

const ubuntuReposContent = `# Sources from Quickonf
deb <mirror> <oscodename> main restricted universe multiverse
deb <mirror> <oscodename>-updates main restricted universe multiverse
deb <mirror> <oscodename>-backports main restricted universe multiverse
deb http://security.ubuntu.com/ubuntu/ <oscodename>-security main restricted universe multiverse
`

var _ = register(
	"ubuntu repos",
	"ubuntu",
	"Enable base Ubuntu repositories",
	ParamsDesc{
		{"mirror", "Ubuntu mirror URL", ParamTypeString},
	},
	func(params map[string]any, vars variables.Variables, msg status.SendStatus) Command {
		mirror := params["mirror"].(string)
		vars.Define("mirror", mirror)
		return &ubuntuRepos{
			msg:    msg,
			mirror: mirror,
			backend: &fileWatcher{
				cmd: &fileContent{
					fileWatcherCmdInit(
						map[string]any{"path": "/etc/apt/sources.list"},
						vars, msg,
					),
					ubuntuReposContent,
					"root",
				},
			},
		}
	},
)

type ubuntuRepos struct {
	msg status.SendStatus

	mirror string

	backend Command
}

func (u *ubuntuRepos) UpdateVariables(vars variables.Variables) {
	vars.Define("mirror", u.mirror)
	u.backend.UpdateVariables(vars)
}

func (u *ubuntuRepos) Watch() {
	u.backend.Watch()
}

func (u *ubuntuRepos) Stop() {
	u.backend.Stop()
}

func (u *ubuntuRepos) Apply() bool {
	if ok := u.backend.Apply(); !ok {
		return false
	}
	return external.AptGet(
		func(s status.Status, info string, detail status.Detail) {
			if info == "" {
				switch s {
				case status.StatusRunning:
					info = "Downloading packages list from mirror " + u.mirror
				case status.StatusApplied:
					info = "Successfully applied sources for mirror " + u.mirror
				case status.StatusFailed:
					info = "Failed downloading packages list from mirror " + u.mirror
				}
			}
			u.msg(s, info, detail, nil)
		},
		"update",
	)
}
