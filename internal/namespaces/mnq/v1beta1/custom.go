package mnq

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mnq").Groups = []string{"integration"}

	human.RegisterMarshalerFunc(
		mnq.SnsInfoStatus(""),
		human.EnumMarshalFunc(mnqSqsInfoStatusMarshalSpecs),
	)

	cmds.MustFind("mnq", "nats", "get-account").Override(mnqNatsGetAccountBuilder)
	cmds.MustFind("mnq", "nats", "list-credentials").Override(mnqNatsListCredentialsBuilder)

	cmds.MustFind("mnq", "sqs", "list-credentials").Override(mnqSqsListCredentialsBuilder)
	cmds.MustFind("mnq", "sqs", "get-credentials").Override(mnqSqsGetCredentialsBuilder)

	cmds.MustFind("mnq", "sns", "list-credentials").Override(mnqSnsListCredentialsBuilder)

	cmds.MustFind("mnq", "sns", "get-credentials").Override(mnqSnsGetCredentialsBuilder)

	cmds.Merge(core.NewCommands(
		createContextCommand(),
	))

	return cmds
}
