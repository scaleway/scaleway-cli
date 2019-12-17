package instance

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func applyCustomRuns(c *core.Commands) {
	c.MustFind("instance", "security-group", "get").Run = customInstanceGetSecurityGroupRun
	c.MustFind("instance", "image", "list").Run = customInstanceListImagesRun
}

type customSecurityGroupResponse struct {
	instance.SecurityGroup

	Rules []*instance.SecurityGroupRule
}

func customInstanceGetSecurityGroupRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	req := argsI.(*instance.GetSecurityGroupRequest)

	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)
	securityGroup, err := api.GetSecurityGroup(req)
	if err != nil {
		return nil, err
	}

	securityGroupRules, err := api.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
		SecurityGroupID: securityGroup.SecurityGroup.ID,
	}, scw.WithAllPages())
	if err != nil {
		return nil, err
	}

	return &customSecurityGroupResponse{
		SecurityGroup: *securityGroup.SecurityGroup,
		Rules:         securityGroupRules.Rules,
	}, nil
}

// customInstanceListImages list the images for a given organization.
// A call to GetServer(..) with the ID contained in Image.FromServer retrieves more information about the server.
func customInstanceListImagesRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	// customImage is based on instance.Image, with additional information about the server
	type customImage struct {
		ID                string
		Name              string
		Arch              instance.Arch
		CreationDate      time.Time
		ModificationDate  time.Time
		DefaultBootscript *instance.Bootscript
		ExtraVolumes      map[string]*instance.Volume
		Organization      string
		Public            bool
		RootVolume        *instance.VolumeSummary
		State             instance.ImageState
		// Replace Image.FromServer
		ServerID   string
		ServerName string
		Zone       scw.Zone
	}

	// Get images
	req := argsI.(*instance.ListImagesRequest)
	req.Public = scw.BoolPtr(false)
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)
	listImagesResponse, err := api.ListImages(req, scw.WithAllPages())
	if err != nil {
		return nil, err
	}
	images := listImagesResponse.Images

	// Builds customImages
	customImages := []*customImage(nil)
	for _, image := range images {
		customImage_ := &customImage{
			ID:                image.ID,
			Name:              image.Name,
			Arch:              image.Arch,
			CreationDate:      image.CreationDate,
			ModificationDate:  image.ModificationDate,
			DefaultBootscript: image.DefaultBootscript,
			ExtraVolumes:      image.ExtraVolumes,
			Organization:      image.Organization,
			Public:            image.Public,
			RootVolume:        image.RootVolume,
			State:             image.State,
		}

		if image.FromServer != "" {
			serverReq := instance.GetServerRequest{
				Zone:     req.Zone,
				ServerID: image.FromServer,
			}
			getServerResponse, err := api.GetServer(&serverReq)
			if err != nil {
				return nil, err
			}
			zone := scw.Zone("")
			if getServerResponse.Server.Location != nil {
				zone_, err := scw.ParseZone(getServerResponse.Server.Location.ZoneID)
				if err != nil {
					return nil, err
				}
				zone = zone_
			}
			customImage_.ServerID = getServerResponse.Server.ID
			customImage_.ServerName = getServerResponse.Server.Name
			customImage_.Zone = zone
		}
		customImages = append(customImages, customImage_)
	}

	return customImages, nil
}
