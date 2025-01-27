package iot_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/iot/v1"
)

func Test_CreateNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   iot.GetCommands(),
		BeforeFunc: createHub(),
		Cmd:        "scw iot network create hub-id={{ .Hub.ID }} type=sigfox topic-prefix=foo",
		Check:      core.TestCheckGolden(),
		AfterFunc:  deleteHub(),
	}))
}
