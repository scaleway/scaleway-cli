package test

import (
	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	return GetGeneratedCommands()
}
