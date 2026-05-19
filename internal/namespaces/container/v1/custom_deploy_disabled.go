//go:build !(darwin || linux || windows)

package container

import "github.com/scaleway/scaleway-cli/v2/core"

func containerDeployCommand() *core.Command {
	return nil
}
