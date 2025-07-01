package redis

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func redisACLUpdateCommand() *core.Command {
	return &core.Command{
		Short:     "Update an ACL rule for a Redis™ Database Instance (network rule)",
		Long:      "Update an ACL rule (IP/description) for a Redis™ Database Instance (Redis™ cluster). This command simulates an update by fetching, deleting, and re-adding the rule.",
		Namespace: "redis",
		Resource:  "acl",
		Verb:      "update",
		ArgsType:  reflect.TypeOf(redisUpdateACLRuleRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      "UUID of the Redis cluster",
				Required:   true,
				Positional: false,
			},
			{
				Name:       "acl-id",
				Short:      "UUID of the ACL rule to update",
				Required:   true,
				Positional: true,
			},
			{
				Name:     "ip-cidr",
				Short:    "New IPv4 network address of the rule (optional, defaults to current)",
				Required: false,
			},
			{
				Name:     "description",
				Short:    "New description of the rule (optional, defaults to current)",
				Required: false,
			},
			{
				Name:     "zone",
				Short:    "Zone to target. If none is passed will use default zone from the config (fr-par-1 | fr-par-2 | nl-ams-1 | nl-ams-2 | pl-waw-1 | pl-waw-2)",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*redisUpdateACLRuleRequest)
			api := redis.NewAPI(core.ExtractClient(ctx))

			// 1. Get the existing rule
			rule, err := api.GetACLRule(&redis.GetACLRuleRequest{
				ACLID: args.ACLID,
				Zone:  scw.Zone(args.Zone),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to get ACL rule: %w", err)
			}

			// 2. Delete the existing rule
			_, err = api.DeleteACLRule(&redis.DeleteACLRuleRequest{
				ACLID: args.ACLID,
				Zone:  scw.Zone(args.Zone),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to delete ACL rule: %w", err)
			}

			waitReq := &redis.WaitForClusterRequest{
				ClusterID: args.ClusterID,
				Zone:      scw.Zone(args.Zone),
				Timeout:   scw.TimeDurationPtr(redisActionTimeout),
			}
			_, err = api.WaitForCluster(waitReq)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for cluster to be ready: %w", err)
			}

			// 3. Add a new rule with updated fields
			ipCIDR := args.IPCidr
			if ipCIDR == "" {
				ipCIDR = rule.IPCidr.String()
			}
			desc := args.Description
			if desc == "" {
				desc = *rule.Description
			}
			_, ipnet, err := net.ParseCIDR(ipCIDR)
			if err != nil {
				return nil, fmt.Errorf("invalid ip-cidr: %w", err)
			}
			scwIPNet := scw.IPNet{IPNet: *ipnet}
			addResp, err := api.AddACLRules(&redis.AddACLRulesRequest{
				ClusterID: args.ClusterID,
				Zone:      scw.Zone(args.Zone),
				ACLRules: []*redis.ACLRuleSpec{
					{
						IPCidr:      scwIPNet,
						Description: desc,
					},
				},
			})
			if err != nil {
				return nil, fmt.Errorf("failed to add updated ACL rule: %w", err)
			}

			_, err = api.WaitForCluster(waitReq)
			if err != nil {
				return nil, fmt.Errorf("failed to wait for cluster to be ready: %w", err)
			}

			return addResp.ACLRules, nil
		},
	}
}

type redisUpdateACLRuleRequest struct {
	ClusterID   string
	ACLID       string
	IPCidr      string
	Description string
	Zone        string
}
