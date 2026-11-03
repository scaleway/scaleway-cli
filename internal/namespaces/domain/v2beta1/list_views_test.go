package domain_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/domain/v2beta1"
)

func Test_ListCertificate(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: domain.GetCommands(),
		Cmd:      "scw dns certificate list example.com",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListContact(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: domain.GetCommands(),
		Cmd:      "scw domain contact list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDomain(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: domain.GetCommands(),
		Cmd:      "scw domain domain list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListHost(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: domain.GetCommands(),
		Cmd:      "scw domain host list domain=example.com",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListTask(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: domain.GetCommands(),
		Cmd:      "scw domain task list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListTld(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: domain.GetCommands(),
		Cmd:      "scw domain tld list",
		Check:    core.TestCheckGolden(),
	}))
}
