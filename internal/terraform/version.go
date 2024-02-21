package terraform

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// GetClientVersion runs terraform version and returns the version string
func GetLocalClientVersion() (*Version, error) {
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
