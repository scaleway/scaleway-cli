package namespaces

import (
	document_db "github.com/scaleway/scaleway-cli/v2/internal/namespaces/document_db/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"os"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	accountv3 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/alias"
	applesilicon "github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
	autocompleteNamespace "github.com/scaleway/scaleway-cli/v2/internal/namespaces/autocomplete"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2alpha1"
	cockpit "github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1beta1"
	configNamespace "github.com/scaleway/scaleway-cli/v2/internal/namespaces/config"
	container "github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1"
	domain "github.com/scaleway/scaleway-cli/v2/internal/namespaces/domain/v2beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/feedback"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
	function "github.com/scaleway/scaleway-cli/v2/internal/namespaces/function/v1beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/help"
	iam "github.com/scaleway/scaleway-cli/v2/internal/namespaces/iam/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/info"
	initNamespace "github.com/scaleway/scaleway-cli/v2/internal/namespaces/init"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/iot/v1"
	ipfs "github.com/scaleway/scaleway-cli/v2/internal/namespaces/ipfs/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/marketplace/v2"
	mnq "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mnq/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/shell"
	tem "github.com/scaleway/scaleway-cli/v2/internal/namespaces/tem/v1alpha1"
	versionNamespace "github.com/scaleway/scaleway-cli/v2/internal/namespaces/version"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpcgw/v1"
	webhosting "github.com/scaleway/scaleway-cli/v2/internal/namespaces/webhosting/v1alpha1"
)

var labs = os.Getenv("SCW_ENABLE_LABS") == "true"

// Enable beta in the code when products are in beta
var beta = os.Getenv(scw.ScwEnableBeta) == "true"

// GetCommands returns a list of all commands in the CLI.
// It is used by both scw and scw-qa.
// We can not put it in `core` package as it would result in a import cycle `core` -> `namespaces/autocomplete` -> `core`.
func GetCommands() *core.Commands {
	// Import all commands available in CLI from various packages.
	// NB: Merge order impacts scw usage sort.
	commands := core.NewCommandsMerge(
		iam.GetCommands(),
		instance.GetCommands(),
		baremetal.GetCommands(),
		cockpit.GetCommands(),
		k8s.GetCommands(),
		marketplace.GetCommands(),
		initNamespace.GetCommands(),
		configNamespace.GetCommands(),
		accountv3.GetCommands(),
		autocompleteNamespace.GetCommands(),
		object.GetCommands(),
		versionNamespace.GetCommands(),
		registry.GetCommands(),
		feedback.GetCommands(),
		info.GetCommands(),
		rdb.GetCommands(),
		lb.GetCommands(),
		iot.GetCommands(),
		help.GetCommands(),
		vpc.GetCommands(),
		domain.GetCommands(),
		applesilicon.GetCommands(),
		flexibleip.GetCommands(),
		container.GetCommands(),
		function.GetCommands(),
		vpcgw.GetCommands(),
		redis.GetCommands(),
		secret.GetCommands(),
		shell.GetCommands(),
		tem.GetCommands(),
		mnq.GetCommands(),
		alias.GetCommands(),
		webhosting.GetCommands(),
		billing.GetCommands(),
	)
	if labs {
		commands.Merge(ipfs.GetCommands())
	}
	if beta {
		commands.Merge(document_db.GetCommands())
	}

	return commands
}
