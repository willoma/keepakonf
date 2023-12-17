package variables

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
)

var (
	globalMap Variables
	Global    []Variable
)

func init() {
	globalMap = Variables{}

	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Could not get hostname", "error", err)
	}
	globalMap.Define("hostname", hostname)

	lsbRelease, err := exec.Command("lsb_release", "--codename", "--short").Output()
	if err != nil {
		slog.Error("Could not get distribution codename", "error", err)
	}
	globalMap.Define("oscodename", string(bytes.TrimSpace(lsbRelease)))

	distro, err := exec.Command("lsb_release", "--id", "--short").Output()
	if err != nil {
		slog.Error("Could not get distribution name", "error", err)
	}
	globalMap.Define("osdistribution", string(bytes.TrimSpace(distro)))

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
	return globalMap.Clone()
}
