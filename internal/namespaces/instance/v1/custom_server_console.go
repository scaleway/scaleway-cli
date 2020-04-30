package instance

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/gotty"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type instanceConsoleServerArgs struct {
	Zone     scw.Zone
	ServerID string
}

func serverConsoleCommand() *core.Command {
	return &core.Command{
		Short:     `Connect to the serial console of an instance`,
		Namespace: "instance",
		Verb:      "console",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(instanceConsoleServerArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      "Server ID to connect to",
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec(),
		},
		Run: instanceServerConsoleRun,
	}
}

func instanceServerConsoleRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*instanceConsoleServerArgs)

	client := core.ExtractClient(ctx)
	apiInstance := instance.NewAPI(client)
	serverResp, err := apiInstance.GetServer(&instance.GetServerRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}
	server := serverResp.Server

	secretKey, ok := client.GetSecretKey()
	if !ok {
		return nil, fmt.Errorf("could not get secret key")
	}

	ttyClient, err := gotty.NewClient(server.Zone, server.ID, secretKey)
	if err != nil {
		return nil, err
	}

	// Add hint on how to quit properly
	fmt.Printf(terminal.Style("Open connection to %s (%s)\n", color.Bold), server.Name, server.ID)
	fmt.Println(" - You may need to hit enter to start")
	fmt.Println(" - Type Ctrl+q to quit.")
	fmt.Println(interactive.Line("-"))

	if err = ttyClient.Connect(); err != nil {
		return nil, err
	}

	return nil, err
}
