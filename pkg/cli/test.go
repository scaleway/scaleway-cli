package cli

import (
	"fmt"
	"os/exec"
	"syscall"
)

var (
	scwcli         = "../../scw"
	publicCommands = []string{
		"help", "attach", "commit", "cp", "create",
		"events", "exec", "history", "images", "info",
		"inspect", "kill", "login", "logout", "logs",
		"port", "products", "ps", "rename", "restart",
		"rm", "rmi", "run", "search", "start", "stop",
		"tag", "top", "version", "wait",
	}
	secretCommands = []string{
		"_patch", "_completion", "_flush-cache", "_userdata", "_billing",
	}
	publicOptions = []string{
		"-h, --help=false",
		"-D, --debug=false",
		"-V, --verbose=false",
		"-q, --quiet=false",
		"--sensitive=false",
		"-v, --version=false",
	}
)

func shouldFitInTerminal(actual interface{}, expected ...interface{}) string {
	if len(actual.(string)) < 80 {
		return ""
	}
	return fmt.Sprintf("len(%q)\n -> %d chars (> 80 chars)", actual, len(actual.(string)))
}

func getExitCode(err error) (int, error) {
	exitCode := 0
	if exiterr, ok := err.(*exec.ExitError); ok {
		if procExit, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return procExit.ExitStatus(), nil
		}
	}
	return exitCode, fmt.Errorf("failed to get exit code")
}

func processExitCode(err error) (exitCode int) {
	if err != nil {
		var exiterr error
		if exitCode, exiterr = getExitCode(err); exiterr != nil {
			// TODO: Fix this so we check the error's text.
			// we've failed to retrieve exit code, so we set it to 127
			exitCode = 127
		}
	}
	return
}

func runCommandWithOutput(cmd *exec.Cmd) (output string, exitCode int, err error) {
	exitCode = 0
	out, err := cmd.CombinedOutput()
	exitCode = processExitCode(err)
	output = string(out)
	return
}
