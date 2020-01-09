package instance

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	serverActionTimeout = 10 * time.Minute
)

func getCustomCommands() *core.Commands {
	return core.NewCommands(
		instanceServerCreate(),
		instanceServerStart(),
		instanceServerStop(),
		instanceServerStandby(),
		instanceServerReboot(),
		instanceSecurityGroupClear(),
		instanceSecurityGroupUpdate(),
		instanceUserData(),
		instanceUserDataList(),
		instanceUserDataSet(),
		instanceUserDataDelete(),
		instanceUserDataGet(),
		instanceServerDelete(),
	)
}

func instanceUserData() *core.Command {
	return &core.Command{
		Namespace: "instance",
		Resource:  "user-data",
	}
}

type instanceActionRequest struct {
	Zone     scw.Zone
	ServerID string
}

var serverActionArgSpecs = core.ArgSpecs{
	{
		Name:     "server-id",
		Short:    `ID of the server affected by the action.`,
		Required: true,
	},
	core.ZoneArgSpec,
}

func instanceServerStart() *core.Command {
	return &core.Command{
		Short:     `Power on server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "start",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionPoweron),
		WaitFunc:  waitForServerFunc,
		ArgSpecs:  serverActionArgSpecs,
	}
}

func instanceServerStop() *core.Command {
	return &core.Command{
		Short:     `Power off server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "stop",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionPoweroff),
		WaitFunc:  waitForServerFunc,
		ArgSpecs:  serverActionArgSpecs,
	}
}

func instanceServerStandby() *core.Command {
	return &core.Command{
		Short:     `Put server in standby mode`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "standby",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionStopInPlace),
		WaitFunc:  waitForServerFunc,
		ArgSpecs:  serverActionArgSpecs,
	}
}

func instanceServerReboot() *core.Command {
	return &core.Command{
		Short:     `Reboot server`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "reboot",
		ArgsType:  reflect.TypeOf(instanceActionRequest{}),
		Run:       getRunServerAction(instance.ServerActionReboot),
		WaitFunc:  waitForServerFunc,
		ArgSpecs:  serverActionArgSpecs,
	}
}

func waitForServerFunc(ctx context.Context, argsI, _ interface{}) error {
	_, err := instance.NewAPI(core.ExtractClient(ctx)).WaitForServer(&instance.WaitForServerRequest{
		Zone:     argsI.(*instanceActionRequest).Zone,
		ServerID: argsI.(*instanceActionRequest).ServerID,
		Timeout:  serverActionTimeout,
	})
	return err
}

func getRunServerAction(action instance.ServerAction) core.CommandRunner {
	return func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		args := argsI.(*instanceActionRequest)

		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)

		_, err := api.ServerAction(&instance.ServerActionRequest{
			Zone:     args.Zone,
			ServerID: args.ServerID,
			Action:   action,
		})
		return &core.SuccessResult{Message: fmt.Sprintf("%s successful for the server", action)}, err
	}
}

type instanceResetSecurityGroupArgs struct {
	Zone            scw.Zone
	SecurityGroupID string
}

func instanceSecurityGroupClear() *core.Command {
	return &core.Command{
		Short:     `Remove all rules of a security group`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "clear",
		ArgsType:  reflect.TypeOf(instanceResetSecurityGroupArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*instanceResetSecurityGroupArgs)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			rules, err := api.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
				Zone:            args.Zone,
				SecurityGroupID: args.SecurityGroupID,
			}, scw.WithAllPages())
			if err != nil {
				return nil, err
			}

			var deleteErrors error
			for _, rule := range rules.Rules {
				if !rule.Editable {
					continue
				}
				err = api.DeleteSecurityGroupRule(&instance.DeleteSecurityGroupRuleRequest{
					Zone:                args.Zone,
					SecurityGroupID:     args.SecurityGroupID,
					SecurityGroupRuleID: rule.ID,
				})
				if err != nil {
					deleteErrors = multierror.Append(deleteErrors, err)
				}
			}
			if deleteErrors != nil {
				return nil, deleteErrors
			}
			return &core.SuccessResult{Message: "Successful reset of the security group rules"}, err
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "security-group-id",
				Short:    `ID of the security group to reset.`,
				Required: true,
			},
		},
	}
}

func instanceSecurityGroupUpdate() *core.Command {
	return &core.Command{
		Short:     `Update security group`,
		Long:      `Update security group.`,
		Namespace: "instance",
		Resource:  "security-group",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(instance.UpdateSecurityGroupRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec,
			{
				Name:     "security-group-id",
				Short:    `ID of the security group to update`,
				Required: true,
			},
			{
				Name: "name",
			},
			{
				Name: "description",
			},
			{
				Name: "stateful",
			},
			{
				Name:       "inbound-default-policy",
				Default:    core.DefaultValueSetter("accept"),
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name:       "outbound-default-policy",
				Default:    core.DefaultValueSetter("accept"),
				EnumValues: []string{"accept", "drop"},
			},
			{
				Name: "organization-default",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			req := argsI.(*instance.UpdateSecurityGroupRequest)

			api := instance.NewAPI(core.ExtractClient(ctx))
			res, err := api.UpdateSecurityGroup(req)
			if err == nil {
				return res, nil
			}

			resErr, isResErr := err.(*scw.ResponseError)
			if !isResErr {
				return nil, err
			}

			// Try to find the error type and create a more user friendly one.
			switch resErr.Message {
			case "default security group can't be stateful":
				return nil, &core.CliError{
					Err: fmt.Errorf("your default security group cannot be stateful"),
					Details: interactive.RemoveIndent(`
						You have to make this security group stateless to use it as an organization default.
						More info: https://www.scaleway.com/en/docs/how-to-activate-a-stateful-cloud-firewall
					`),
					Hint: "scw instance security-group update security-group-id=" + req.SecurityGroupID + " organization-default=true stateful=false",
				}

			case "cannot have more than one organization default":
				defaultSG, err := getDefaultOrganizationSecurityGroup(ctx, req.Zone)
				if err != nil {
					// Abort and return the original error.
					return nil, resErr
				}

				return nil, &core.CliError{
					Err: fmt.Errorf("you cannot have more than one organization default"),
					Details: interactive.RemoveIndent(`
						You already have an organization default security-group (` + defaultSG.ID + `).

						First, you need to set your current organization default security-group as non-default with:
						scw instance security-group update security-group-id=` + defaultSG.ID + ` organization-default=false

						Then, retry this command:
						scw instance security-group update security-group-id=` + req.SecurityGroupID + ` organization-default=true stateful=false
					`),
				}
			default:
				// Unknown error, use default behavior.
				return nil, resErr
			}
		},
	}
}

func getDefaultOrganizationSecurityGroup(ctx context.Context, zone scw.Zone) (*instance.SecurityGroup, error) {
	api := instance.NewAPI(core.ExtractClient(ctx))

	orgID := core.GetOrganizationIdFromContext(ctx)
	sgList, err := api.ListSecurityGroups(&instance.ListSecurityGroupsRequest{
		Zone:         zone,
		Organization: scw.StringPtr(orgID),
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	for _, sg := range sgList.SecurityGroups {
		if sg.OrganizationDefault {
			return sg, nil
		}
	}

	return nil, fmt.Errorf("%s organization does not have a default security group", orgID)
}

func instanceUserDataList() *core.Command {
	return &core.Command{
		Short:     `List user data`,
		Long:      `List user data for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(instance.ListServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec,
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return instance.NewAPI(core.ExtractClient(ctx)).ListServerUserData(argsI.(*instance.ListServerUserDataRequest))
		},
	}
}

func instanceUserDataDelete() *core.Command {
	return &core.Command{
		Short:     `Delete user data by key`,
		Long:      `Delete user data key for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "delete",
		ArgsType:  reflect.TypeOf(instance.DeleteServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec,
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
			{
				Name:     "key",
				Short:    `Key of the user data`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			err := instance.NewAPI(core.ExtractClient(ctx)).DeleteServerUserData(argsI.(*instance.DeleteServerUserDataRequest))
			if err != nil {
				return nil, err
			}
			return &core.SuccessResult{}, nil
		},
	}
}

func instanceUserDataGet() *core.Command {
	return &core.Command{
		Short:     `Get user data key`,
		Long:      `Get user data key for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(instance.GetServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec,
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
			{
				Name:     "key",
				Short:    `Key of the user data`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return instance.NewAPI(core.ExtractClient(ctx)).GetServerUserData(argsI.(*instance.GetServerUserDataRequest))
		},
	}
}

func instanceUserDataSet() *core.Command {
	return &core.Command{
		Short:     `Set a user data`,
		Long:      `Set a user data for the given server.`,
		Namespace: "instance",
		Resource:  "user-data",
		Verb:      "set",
		ArgsType:  reflect.TypeOf(instance.SetServerUserDataRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec,
			{
				Name:     "server-id",
				Short:    `ID of a server`,
				Required: true,
			},
			{
				Name:     "key",
				Short:    `Key of the user data`,
				Required: true,
			},
			{
				Name:     "content",
				Short:    `Content of the user data`,
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			err := instance.NewAPI(core.ExtractClient(ctx)).SetServerUserData(argsI.(*instance.SetServerUserDataRequest))
			if err != nil {
				return nil, err
			}
			return &core.SuccessResult{}, nil
		},
	}
}

type customeDeleteServerRequest struct {
	Zone          scw.Zone
	ServerID      string
	DeleteIP      bool
	DeleteVolumes bool
}

func instanceServerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete server`,
		Long:      `Delete a server with the given ID.`,
		Namespace: "instance",
		Verb:      "delete",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(customeDeleteServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec,
			{
				Name:     "server-id",
				Required: true,
			},
			{
				Name:  "delete-ip",
				Short: "Delete the IP attached to the server as well",
			},
			{
				Name:  "delete-volumes",
				Short: "Delete the volumes attached to the server as well",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw instance server stop",
				Short:   "Stop a running server",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*customeDeleteServerRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			server, err := api.GetServer(&instance.GetServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, err
			}

			err = api.DeleteServer(&instance.DeleteServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			if err != nil {
				return nil, err
			}

			var multiErr error
			if args.DeleteIP && server.Server.PublicIP != nil {
				err = api.DeleteIP(&instance.DeleteIPRequest{
					Zone: args.Zone,
					IP:   server.Server.PublicIP.ID,
				})
				if err != nil {
					multiErr = multierror.Append(multiErr, err)
				}
			}

			if args.DeleteVolumes {
				for _, volume := range server.Server.Volumes {
					err = api.DeleteVolume(&instance.DeleteVolumeRequest{
						Zone:     args.Zone,
						VolumeID: volume.ID,
					})
					if err != nil {
						multiErr = multierror.Append(multiErr, err)
					}
				}
			}
			if multiErr != nil {
				return nil, &core.CliError{
					Err:  multiErr,
					Hint: "Make sure these resources have been deleted or try to delete it manually.",
				}
			}

			return &core.SuccessResult{}, nil
		},
	}
}
