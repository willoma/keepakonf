package variables

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
)

var (
	globalMap map[string]string
	Global    []Variable
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Could not get hostname", "error", err)
	}

	lsbRelease, err := exec.Command("lsb_release", "--codename", "--short").Output()
	if err != nil {
		slog.Error("Could not get distribution codename", "error", err)
	}

	distro, err := exec.Command("lsb_release", "--id", "--short").Output()
	if err != nil {
		slog.Error("Could not get distribution name", "error", err)
	}

	globalMap = map[string]string{
		"hostname":       hostname,
		"oscodename":     string(bytes.TrimSpace(lsbRelease)),
		"osdistribution": string(bytes.TrimSpace(distro)),
	}
	Global = []Variable{
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

func GlobalMap() Variables {
	gvm := make(Variables, len(globalMap))
	for k, v := range globalMap {
		gvm.Define(k, v)
	}
	return gvm
}
