package vpcgw

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var vpcgwPATRulesEditYamlExample = `pat_rules:
- public_port: 2222
  private_ip: 192.168.1.1
  private_port: 22
  protocol: tcp
`

type vpcgwPATRulesEditArgs struct {
	Zone      scw.Zone
	GatewayID string
	Mode      editor.MarshalMode
}

func vpcgwPATRulesEditCommand() *core.Command {
	return &core.Command{
		Short:     "Edit all PAT rules of a Public Gateway",
		Long:      editor.LongDescription,
		Namespace: "vpc-gw",
		Resource:  "pat-rule",
		Verb:      "edit",
		ArgsType:  reflect.TypeOf(vpcgwPATRulesEditArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "gateway-id",
				Short:      "ID of the PAT rules' Public Gateway",
				Required:   true,
				Positional: true,
			},
			editor.MarshalModeArgSpec(),
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*vpcgwPATRulesEditArgs)

			client := core.ExtractClient(ctx)
			api := vpcgw.NewAPI(client)

			setRequest := &vpcgw.SetPatRulesRequest{
				Zone:      args.Zone,
				GatewayID: args.GatewayID,
			}

			rules, err := api.ListPatRules(&vpcgw.ListPatRulesRequest{
				Zone:       args.Zone,
				GatewayIDs: []string{args.GatewayID},
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to list PAT rules: %w", err)
			}

			editedSetRequest, err := editor.UpdateResourceEditor(rules, setRequest, &editor.Config{
				PutRequest:  true,
				MarshalMode: args.Mode,
				Template:    vpcgwPATRulesEditYamlExample,
			})
			if err != nil {
				return nil, err
			}

			setRequest = editedSetRequest.(*vpcgw.SetPatRulesRequest)

			resp, err := api.SetPatRules(setRequest, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to set PAT rules: %w", err)
			}

			return resp.PatRules, nil
		},
	}
}
