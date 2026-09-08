package instance

import (
	"context"
	"net"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

var privateNICStateMarshalSpecs = human.EnumMarshalSpecs{
	instance.PrivateNICStateAvailable:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.PrivateNICStateSyncing:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.PrivateNICStateSyncingError: &human.EnumMarshalSpec{Attribute: color.FgRed},
}

func privateNicGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("private-nic-id").Short = "The private NIC unique ID or MAC address"

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		tmpRequest := argsI.(*instance.GetPrivateNICRequest)

		if isMacAddress(tmpRequest.PrivateNicID) {
			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			listPrivateNICs, err := api.ListPrivateNICs(&instance.ListPrivateNICsRequest{
				Zone:     tmpRequest.Zone,
				ServerID: tmpRequest.ServerID,
			})
			if err != nil {
				return nil, err
			}
			for _, pn := range listPrivateNICs.PrivateNics {
				if pn.MacAddress == tmpRequest.PrivateNicID {
					tmpRequest.PrivateNicID = pn.ID
				}
			}
		}

		return runner(ctx, tmpRequest)
	}

	return deprecatePrivateNicCommand(c)
}

func isMacAddress(address string) bool {
	_, err := net.ParseMAC(address)

	return err == nil
}

func privateNicBuilder(c *core.Command) *core.Command {
	return deprecatePrivateNicCommand(c)
}

func privateNicCreateBuilder(c *core.Command) *core.Command {
	return deprecatePrivateNicCommand(c)
}

func privateNicListBuilder(c *core.Command) *core.Command {
	return deprecatePrivateNicCommand(c)
}

func privateNicUpdateBuilder(c *core.Command) *core.Command {
	return deprecatePrivateNicCommand(c)
}

func privateNicDeleteBuilder(c *core.Command) *core.Command {
	return deprecatePrivateNicCommand(c)
}

func deprecatePrivateNicCommand(c *core.Command) *core.Command {
	c.Hidden = true
	c.Deprecated = true
	c.DeprecationMessage = "private-nic is a term only used in v1, please use private-network-interface instead."

	return c
}
