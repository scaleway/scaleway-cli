package iot

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_CreateNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createHub(),
		Cmd:        "scw iot network create hub-id={{ .Hub.ID }} type=sigfox topic-prefix=foo",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteHub(),
	}))
}
