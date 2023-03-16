package editor_test

import (
	"context"
	"log"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

func TestInterceptor(t *testing.T) {
	cmds := namespaces.GetCommands()
	interceptor := editor.Interceptor(cmds.MustFind("container", "namespace", "get"))
	ctx := context.Background()
	res, err := interceptor(ctx, &container.UpdateNamespaceRequest{
		Region:      "fr-par",
		NamespaceID: "47c8f0ad-3567-451e-a0a5-820318678ce3",
	}, cmds.MustFind("container", "namespace", "update").Run)
	assert.Nil(t, err)
	log.Println(res)
}
