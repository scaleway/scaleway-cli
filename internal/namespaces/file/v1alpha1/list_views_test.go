package file_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/file/v1alpha1"
)

func Test_ListAttachment(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: file.GetCommands(),
		Cmd:      "scw file attachment list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListFilesystem(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: file.GetCommands(),
		Cmd:      "scw file filesystem list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListFilesystemType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: file.GetCommands(),
		Cmd:      "scw file filesystem-type list",
		Check:    core.TestCheckGolden(),
	}))
}
