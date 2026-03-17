package edge_services

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	edgeservices "github.com/scaleway/scaleway-sdk-go/api/edge_services/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var edgeServicesRouteRulesEditYamlExample = `route_rules:
- rule_http_match:
    method_filters:
    - get
    - post
    path_filter:
      path_filter_type: regex
      value: ^/api/.*
  backend_stage_id: 11111111-1111-1111-1111-111111111111
- rule_http_match:
    method_filters:
    - get
    path_filter:
      path_filter_type: regex
      value: ^/static/.*
  backend_stage_id: 11111111-1111-1111-1111-111111111111
`

var edgeServicesRouteRulesEditYamlExampleSimple = `route_rules:
- rule_http_match:
    method_filters:
    - get
    - post
    path_filter:
      path_filter_type: regex
      value: ^/api/.*
- rule_http_match:
    method_filters:
    - get
    path_filter:
      path_filter_type: regex
      value: ^/static/.*
`

type edgeServicesRouteRulesEditArgs struct {
	RouteStageID   string
	BackendStageID *string
	Mode           editor.MarshalMode
}

func edgeServicesRouteRulesEditCommand() *core.Command {
	return &core.Command{
		Short: "Edit all route rules of a route stage",
		Long: `Edit all route rules of a route stage.

If backend-stage-id is provided, the editor will only show rule_http_match fields and the specified backend will be applied to all rules automatically.
Otherwise, opens the editor with full rules including backend_stage_id for each rule.`,
		Namespace: "edge-services",
		Resource:  "route-rules",
		Verb:      "edit",
		ArgsType:  reflect.TypeOf(edgeServicesRouteRulesEditArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "route-stage-id",
				Short:      "ID of the route stage to edit",
				Required:   true,
				Positional: true,
			},
			{
				Name:     "backend-stage-id",
				Short:    "ID of the backend stage to apply to all rules (simplifies editing when using a single backend)",
				Required: false,
			},
			editor.MarshalModeArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			args := argsI.(*edgeServicesRouteRulesEditArgs)

			client := core.ExtractClient(ctx)
			api := edgeservices.NewAPI(client)

			rules, err := api.ListRouteRules(&edgeservices.ListRouteRulesRequest{
				RouteStageID: args.RouteStageID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to list route rules: %w", err)
			}

			setRequest := &edgeservices.SetRouteRulesRequest{
				RouteStageID: args.RouteStageID,
			}

			template := edgeServicesRouteRulesEditYamlExample
			var ignoreFields []string
			if args.BackendStageID != nil {
				template = edgeServicesRouteRulesEditYamlExampleSimple
				ignoreFields = []string{"backend_stage_id"}
			}

			editedSetRequest, err := editor.UpdateResourceEditor(rules, setRequest, &editor.Config{
				PutRequest:   true,
				MarshalMode:  args.Mode,
				Template:     template,
				IgnoreFields: ignoreFields,
			})
			if err != nil {
				return nil, err
			}

			setRequest = editedSetRequest.(*edgeservices.SetRouteRulesRequest)

			if args.BackendStageID != nil {
				for _, rule := range setRequest.RouteRules {
					rule.BackendStageID = args.BackendStageID
				}
			}

			resp, err := api.SetRouteRules(setRequest, scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("failed to set route rules: %w", err)
			}

			return resp.RouteRules, nil
		},
	}
}
