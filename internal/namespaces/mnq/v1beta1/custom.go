package mnq

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(mnq.SnsInfoStatus(""), human.EnumMarshalFunc(mnqSqsInfoStatusMarshalSpecs))

	cmds.Merge(core.NewCommands(
		aliasCreateContextCommand(),
	))

	return cmds
}
