package mongodb

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/passwordgenerator"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	instanceActionTimeout = 20 * time.Minute
)

var instanceStatusMarshalSpecs = human.EnumMarshalSpecs{
	mongodb.InstanceStatusConfiguring: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "configuring",
	},
	mongodb.InstanceStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "deleting",
	},
	mongodb.InstanceStatusError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "error",
	},
	mongodb.InstanceStatusInitializing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "initializing",
	},
	mongodb.InstanceStatusLocked: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "locked",
	},
	mongodb.InstanceStatusProvisioning: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "provisioning",
	},
	mongodb.InstanceStatusReady: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "ready",
	},
	mongodb.InstanceStatusSnapshotting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "snapshotting",
	},
}

type mongodbCreateInstanceRequestCustom struct {
	*mongodb.CreateInstanceRequest
	GeneratePassword bool
}

type CreateInstanceResult struct {
	*mongodb.Instance
	Password string `json:"password,omitempty"`
}

func instanceCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-amount").Default = core.DefaultValueSetter("1")
	c.ArgSpecs.GetByName("version").Default = fetchLatestEngine
	c.ArgSpecs.GetByName("volume.size-bytes").Default = core.DefaultValueSetter("5GB")
	c.ArgSpecs.GetByName("volume.type").Default = core.DefaultValueSetter("sbs_5k")
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteNodeType
	c.ArgSpecs.AddBefore("password", passwordgenerator.ArgSpec(true))
	c.ArgSpecs.GetByName("password").Required = false
	c.ArgsType = reflect.TypeFor[mongodbCreateInstanceRequestCustom]()

	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		prev := respI.(CreateInstanceResult)
		api := mongodb.NewAPI(core.ExtractClient(ctx))

		instance, err := api.WaitForInstance(&mongodb.WaitForInstanceRequest{
			InstanceID:    prev.ID,
			Region:        prev.Region,
			Timeout:       new(instanceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			return nil, err
		}

		return CreateInstanceResult{Instance: instance, Password: prev.Password}, nil
	}

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		api := mongodb.NewAPI(core.ExtractClient(ctx))
		customRequest := argsI.(*mongodbCreateInstanceRequestCustom)
		req := customRequest.CreateInstanceRequest

		generated := customRequest.GeneratePassword && req.Password == ""
		var err error
		req.Password, err = passwordgenerator.Resolve(customRequest.GeneratePassword, req.Password)
		if err != nil {
			return nil, err
		}

		instance, err := api.CreateInstance(req, scw.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		result := CreateInstanceResult{Instance: instance}
		if generated {
			result.Password = req.Password
		}

		return result, nil
	}

	return c
}

func fetchLatestEngine(ctx context.Context) (string, string) {
	client := core.ExtractClient(ctx)
	api := mongodb.NewAPI(client)
	latestValueVersion, err := api.FetchLatestEngineVersion()
	if err != nil {
		return "", ""
	}

	return latestValueVersion.Version, latestValueVersion.Version
}

var completeListNodeTypeCache *mongodb.ListNodeTypesResponse

func autoCompleteNodeType(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	region := scw.Region("")
	switch req := request.(type) {
	case *mongodb.CreateInstanceRequest:
		region = req.Region
	case *mongodbCreateInstanceRequestCustom:
		region = req.Region
	case *mongodb.UpgradeInstanceRequest:
		region = req.Region
	}

	suggestions := core.AutocompleteSuggestions{}

	client := core.ExtractClient(ctx)
	api := mongodb.NewAPI(client)

	if completeListNodeTypeCache == nil {
		res, err := api.ListNodeTypes(&mongodb.ListNodeTypesRequest{
			Region: region,
		}, scw.WithAllPages())
		if err != nil {
			return nil
		}
		completeListNodeTypeCache = res
	}

	for _, nodeType := range completeListNodeTypeCache.NodeTypes {
		if strings.HasPrefix(nodeType.Name, prefix) {
			suggestions = append(suggestions, core.Suggestion{Value: nodeType.Name})
		}
	}

	return suggestions
}

type serverWaitRequest struct {
	InstanceID string
	Region     scw.Region
	Timeout    time.Duration
}

func instanceWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for an instance to reach a stable state`,
		Long:      `Wait for an instance to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "mongodb",
		Resource:  "instance",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeFor[serverWaitRequest](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := mongodb.NewAPI(core.ExtractClient(ctx))

			return api.WaitForInstance(&mongodb.WaitForInstanceRequest{
				Region:        argsI.(*serverWaitRequest).Region,
				InstanceID:    argsI.(*serverWaitRequest).InstanceID,
				Timeout:       new(argsI.(*serverWaitRequest).Timeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `ID of the instance you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
			core.WaitTimeoutArgSpec(instanceActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for an instance to reach a stable state",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
