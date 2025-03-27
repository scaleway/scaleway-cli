package mongodb

import (
	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
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

func instanceCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("node-number").Default = core.DefaultValueSetter("1")
	c.ArgSpecs.GetByName("version").Default = fetchLatestEngine
	c.ArgSpecs.GetByName("volume.volume-size").Default = core.DefaultValueSetter("5GB")
	c.ArgSpecs.GetByName("volume.volume-type").Default = core.DefaultValueSetter("sbs_5k")
	c.ArgSpecs.GetByName("node-type").AutoCompleteFunc = autoCompleteNodeType

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
	case *mongodb.UpgradeInstanceRequest:
		region = req.Region
	}

	suggestions := core.AutocompleteSuggestions(nil)

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
			suggestions = append(suggestions, nodeType.Name)
		}
	}

	return suggestions
}
