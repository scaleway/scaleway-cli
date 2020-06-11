package baremetal

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	account "github.com/scaleway/scaleway-sdk-go/api/account/v2alpha1"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func serverInstallBuilder(c *core.Command) *core.Command {
	type baremetalInstallServerRequestCustom struct {
		baremetal.InstallServerRequest
		AllSSHKeys *bool
	}

	c.ArgsType = reflect.TypeOf(baremetalInstallServerRequestCustom{})

	c.ArgSpecs.AddBefore("ssh-key-ids.{index}", &core.ArgSpec{
		Name:       "all-ssh-keys",
		Short:      "Add all SSH keys on your baremetal instance (cannot be used with ssh-key-ids)",
		OneOfGroup: "ssh",
	})

	c.ArgSpecs.GetByName("ssh-key-ids.{index}").OneOfGroup = "ssh"
	c.ArgSpecs.GetByName("ssh-key-ids.{index}").Short = "SSH key IDs authorized on the server (cannot be used with all-ssh-keys)"

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		tmpRequest := argsI.(*baremetalInstallServerRequestCustom)

		// SSH keys management
		if tmpRequest.AllSSHKeys != nil && *tmpRequest.AllSSHKeys {
			client := core.ExtractClient(ctx)
			accountapi := account.NewAPI(client)
			orgID, _ := client.GetDefaultOrganizationID()
			listKeys, err := accountapi.ListSSHKeys(&account.ListSSHKeysRequest{
				OrganizationID: &orgID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			var keyIDs []string
			for _, key := range listKeys.SSHKeys {
				keyIDs = append(keyIDs, key.ID)
			}
			tmpRequest.SSHKeyIDs = keyIDs
		}

		return runner(ctx, &tmpRequest.InstallServerRequest)
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))
		return api.WaitForServerInstall(&baremetal.WaitForServerInstallRequest{
			Zone:          argsI.(*baremetalInstallServerRequestCustom).Zone,
			ServerID:      respI.(*baremetal.Server).ID,
			Timeout:       scw.TimeDurationPtr(serverActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
