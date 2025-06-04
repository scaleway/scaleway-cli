package secret

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
)

type customAccessSecretVersionRequest struct {
	secret.AccessSecretVersionRequest
	Field *string
	Raw   bool
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("secret").Groups = []string{"security"}

	cmds.MustFind("secret", "version", "create").Override(secretVersionCreateBuilder)
	cmds.MustFind("secret", "version", "access").Override(secretVersionAccessBuilder)

	return cmds
}
