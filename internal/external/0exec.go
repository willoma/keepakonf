package external

import (
	"os"
	"os/exec"
	"strings"

	"github.com/willoma/keepakonf/internal/status"
)

// func execAs(usr *user.User, env []string, output io.Writer, cmd string, args ...string) error {
// 	args = append([]string{"-u", usr.Username, "--", cmd}, args...)
// 	return execCmd(env, output, "runuser", args...)
// }

func execToMessage(
	runningInfo, successInfo, failureInfo string,
	send status.SendStatus,
	env []string, cmd string, args ...string,
) bool {
	var cmdline strings.Builder
	cmdline.WriteString("root# ")
	if len(env) > 0 {
		for _, env := range env {
			cmdline.WriteString(env)
			cmdline.WriteByte(' ')
		}
	}
	cmdline.WriteString(cmd)
	cmdline.WriteByte(' ')
	for _, arg := range args {
		if strings.Contains(arg, `"`) {
			cmdline.WriteByte('"')
			cmdline.WriteString(arg)
			cmdline.WriteByte('"')
		} else {
			cmdline.WriteString(arg)
		}
		cmdline.WriteByte(' ')
	}
	cmdline.WriteByte('\n')

	c := exec.Command(cmd, args...)
	c.Env = append(os.Environ(), "LANG=C.UTF-8")
	c.Env = append(c.Env, env...)
	w := newSendWriter(runningInfo, send, cmdline.String())
	c.Stdout = w
	c.Stderr = w
	if err := c.Run(); err != nil {
		send(
			status.StatusFailed,
			failureInfo,
			&status.Terminal{
				Command: cmdline.String(),
				Output:  w.Result(),
			},
		)
		return false
	}
	send(
		status.StatusApplied,
		successInfo,
		&status.Terminal{
			Command: cmdline.String(),
			Output:  w.Result(),
		},
	)
	return true
}

// func ExecErr(err error) string {
// 	if err == nil {
// 		return ""
// 	}
// 	var exitErr *exec.ExitError
// 	if errors.As(err, &exitErr) {
// 		return string(exitErr.Stderr)
// 	}
// 	return err.Error()
// }

// func ExecErrCode(err error) int {
// 	if err == nil {
// 		return 0
// 	}
// 	var exitErr *exec.ExitError
// 	if errors.As(err, &exitErr) {
// 		return exitErr.ExitCode()
// 	}
// 	return -1
// }
