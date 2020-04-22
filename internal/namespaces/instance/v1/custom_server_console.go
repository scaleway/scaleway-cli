package instance

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/console"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	ttyURL = "https://tty.scaleway.com/v2"
)

type instanceConsoleServerRequest struct {
	Zone     scw.Zone
	ServerID string
}

func serverConsoleCommand() *core.Command {
	return &core.Command{
		Short:     `Connect to the serial console of an instance`,
		Namespace: "instance",
		Verb:      "console",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(instanceConsoleServerRequest{}),
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
	args := argsI.(*instanceConsoleServerRequest)

	client := core.ExtractClient(ctx)
	apiInstance := instance.NewAPI(client)
	serverResp, err := apiInstance.GetServer(&instance.GetServerRequest{
		Zone:     args.Zone,
		ServerID: args.ServerID,
	})
	if err != nil {
		return nil, err
	}

	secretKey, ok := client.GetSecretKey()
	if !ok {
		return nil, fmt.Errorf("could not get secret key")
	}

	url := fmt.Sprintf("%s/?arg=%s&arg=%s", ttyURL, secretKey, serverResp.Server.ID)
	gottycli, err := console.NewClient(url)
	if err != nil {
		return nil, err
	}

	if err = gottycli.Connect(); err != nil {
		return nil, err
	}

	done := make(chan bool)

	fmt.Println("You are connected, type 'Ctrl+q' to quit.")

	go func() {
		err = gottycli.Loop()
		gottycli.Close()
		done <- true
	}()

	<-done

	return nil, err
}
