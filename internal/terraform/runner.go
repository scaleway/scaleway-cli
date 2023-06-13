package terraform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type RunResponse struct {
	Stdout   string `js:"stdout"`
	Stderr   string `js:"stderr"`
	ExitCode int    `js:"exitCode"`
}

func runCommandInDir(dir string, command string, args ...string) (*RunResponse, error) {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return &RunResponse{
				Stdout:   outb.String(),
				Stderr:   errb.String(),
				ExitCode: exitError.ExitCode(),
			}, nil
		}

		return nil, err
	}

	return &RunResponse{
		Stdout:   outb.String(),
		Stderr:   errb.String(),
		ExitCode: 0,
	}, nil
}

func runTerraformCommand(dir string, args ...string) (*RunResponse, error) {
	return runCommandInDir(dir, "terraform", args...)
}

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// GetClientVersion runs terraform version and returns the version string
func GetVersion() (*Version, error) {
	response, err := runTerraformCommand("", "version", "-json")
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(response.Stdout), &data)
	if err != nil {
		return nil, err
	}

	rawVersion, ok := data["terraform_version"]
	if !ok {
		return nil, errors.New("failed to get terraform version: terraform_version not found")
	}

	version, ok := rawVersion.(string)
	if !ok {
		return nil, errors.New("failed to get terraform version: terraform_version is not a string")
	}

	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("failed to get terraform version: invalid version format '%s'", version)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform version: invalid major version '%s'", parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform version: invalid minor version '%s'", parts[1])
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("failed to get terraform version: invalid patch version '%s'", parts[2])
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func Init(dir string) (*RunResponse, error) {
	return runTerraformCommand(dir, "init")
}

func GenerateConfig(dir string, target string) (*RunResponse, error) {
	return runTerraformCommand(dir, "plan", fmt.Sprintf("-generate-config-out=%s", target))
}
