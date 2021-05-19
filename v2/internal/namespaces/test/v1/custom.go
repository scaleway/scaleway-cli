package test

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	return GetGeneratedCommands()
}
