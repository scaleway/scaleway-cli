package transactional_email

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	transactional_email "github.com/scaleway/scaleway-sdk-go/api/transactional_email/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(transactional_email.DomainStatus(""), human.EnumMarshalFunc(domainStatusMarshalSpecs))
	human.RegisterMarshalerFunc(transactional_email.EmailStatus(""), human.EnumMarshalFunc(emailStatusMarshalSpecs))

	cmds.MustFind("tem", "domain", "get").Override(domainGetBuilder)

	return cmds
}
