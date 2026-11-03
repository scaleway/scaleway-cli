package tem_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/tem/v1alpha1"
)

func Test_ListBlocklists(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: tem.GetCommands(),
		Cmd:      "scw tem blocklists list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDomain(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: tem.GetCommands(),
		Cmd:      "scw tem domain list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListEmail(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: tem.GetCommands(),
		Cmd:      "scw tem email list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListOfferSubscriptions(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: tem.GetCommands(),
		Cmd:      "scw tem offer-subscriptions list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPools(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: tem.GetCommands(),
		Cmd:      "scw tem pools list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListWebhook(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: tem.GetCommands(),
		Cmd:      "scw tem webhook list",
		Check:    core.TestCheckGolden(),
	}))
}
