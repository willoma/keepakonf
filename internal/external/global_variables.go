package external

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/willoma/keepakonf/internal/log"
)

type Variable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

var (
	globalVarsMap   map[string]Variable
	GlobalVariables []Variable
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err, "Could not get hostname")
	}

	lsbRelease, err := exec.Command("lsb_release", "--codename", "--short").Output()
	if err != nil {
		log.Error(err, "Could not get distribution codename")
	}

	distro, err := exec.Command("lsb_release", "--id", "--short").Output()
	if err != nil {
		log.Error(err, "Could not get distribution name")
	}

	globalVarsMap = map[string]Variable{
		"hostname": {
			Name:        "hostname",
			Description: "Name of the computer",
			Value:       hostname,
		},
		"oscodename": {
			Name:        "oscodename",
			Description: "OS codename",
			Value:       string(bytes.TrimSpace(lsbRelease)),
		},
		"osdistribution": {
			Name:        "osdistribution",
			Description: "OS distribution",
			Value:       string(bytes.TrimSpace(distro)),
		},
	}
	GlobalVariables = []Variable{
		{
			Name:        "hostname",
			Description: "Name of the computer",
			Value:       hostname,
		},
		{
			Name:        "oscodename",
			Description: "OS codename",
			Value:       string(bytes.TrimSpace(lsbRelease)),
		},
		{
			Name:        "osdistribution",
			Description: "OS distribution",
			Value:       string(bytes.TrimSpace(distro)),
		},
	}

}
