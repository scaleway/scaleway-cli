package mnq

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mnq", "credential", "create").Override(credentialCreateBuilder)
	cmds.MustFind("mnq", "credential", "get").Override(credentialGetBuilder)

	cmds.MustFind("mnq").Short = "Messaging and Queuing Alpha APIs"

	for _, cmd := range cmds.GetAll() {
		cmd.Namespace = "mnq-v1alpha1"
	}

	return cmds
}
