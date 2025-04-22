package vpc

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var vpcACLEditYamlExample = `default_policy: drop
is_ipv6: false
rules: [
- protocol: TCP
  src_port_low: 0
  src_port_high: 0
  dst_port_low: 80
  dst_port_high: 80
  source: 0.0.0.0/0
  destination: 0.0.0.0/0
  description: Allow HTTP traffic from any source
  action: accept
- protocol: TCP
  src_port_low: 0
  src_port_high: 0
  dst_port_low: 443
  dst_port_high: 443
  source: 0.0.0.0/0
  destination: 0.0.0.0/0
  description: Allow HTTPS traffic from any source
  action: accept
]
`

type vpcACLEditArgs struct {
	Region        scw.Region
	VpcID         string
	IsIPv6        bool
	DefaultPolicy vpc.Action
	Mode          editor.MarshalMode
}

func vpcACLEditCommand() *core.Command {
	return &core.Command{
		Short:     "Edit all ACL rules of a VPC",
		Long:      editor.LongDescription,
		Namespace: "vpc",
		Resource:  "rule",
		Verb:      "edit",
		ArgsType:  reflect.TypeOf(vpcACLEditArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "vpc-id",
				Short:      "ID of the Network ACL's VPC",
				Required:   true,
				Positional: false,
			},
			{
				Name:       "is-ipv6",
				Short:      "Defines whether this set of ACL rules is for IPv6 (false = IPv4). Each Network ACL can have rules for only one IP type",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "default-policy",
				Short:      "Action to take for packets which do not match any rules",
				Required:   false,
				Positional: false,
			},
			editor.MarshalModeArgSpec(),
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*vpcACLEditArgs)

			client := core.ExtractClient(ctx)
			api := vpc.NewAPI(client)

			setRequest := &vpc.SetACLRequest{
				Region:        args.Region,
				VpcID:         args.VpcID,
				IsIPv6:        args.IsIPv6,
				DefaultPolicy: args.DefaultPolicy,
			}

			rules, err := api.GetACL(&vpc.GetACLRequest{
				Region: args.Region,
				VpcID:  args.VpcID,
				IsIPv6: args.IsIPv6,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to list ACL rules: %w", err)
			}

			editedSetRequest, err := editor.UpdateResourceEditor(rules, setRequest, &editor.Config{
				PutRequest:  true,
				MarshalMode: args.Mode,
				Template:    vpcACLEditYamlExample,
			})
			if err != nil {
				return nil, err
			}

			setRequest = editedSetRequest.(*vpc.SetACLRequest)

			resp, err := api.SetACL(setRequest, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to set ACL rules: %w", err)
			}

			return resp.Rules, nil
		},
	}
}
