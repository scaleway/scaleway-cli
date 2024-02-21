package terraform

import (
	"bytes"
	"fmt"
	"os/exec"
)

type runCommandResponse struct {
	Stdout   string `js:"stdout"`
	Stderr   string `js:"stderr"`
	ExitCode int    `js:"exitCode"`
}

func runCommandInDir(dir string, command string, args ...string) (*runCommandResponse, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return &runCommandResponse{
				Stdout:   outb.String(),
				Stderr:   errb.String(),
				ExitCode: exitError.ExitCode(),
			}, nil
		}

		return nil, err
	}

	return &runCommandResponse{
		Stdout:   outb.String(),
		Stderr:   errb.String(),
		ExitCode: 0,
	}, nil
}

func runTerraformCommand(dir string, args ...string) (*runCommandResponse, error) {
	return runCommandInDir(dir, "terraform", args...)
}

func runInitCommand(dir string) (*runCommandResponse, error) {
	return runTerraformCommand(dir, "init")
}

func runGenerateConfigCommand(dir string, target string) (*runCommandResponse, error) {
	return runTerraformCommand(dir, "plan", fmt.Sprintf("-generate-config-out=%s", target))
}
