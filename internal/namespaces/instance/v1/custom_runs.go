package instance

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func applyCustomRuns(c *core.Commands) {
	// Security Group
	c.MustFind("instance", "security-group", "get").Run = customInstanceSecurityGroupGetRun
	instanceSecurityGroupDeleteCmd := c.MustFind("instance", "security-group", "delete")
	instanceSecurityGroupDeleteCmd.Run = customInstanceSecurityGroupDeleteRun(instanceSecurityGroupDeleteCmd.Run)

	// Image
	c.MustFind("instance", "image", "list").Run = customInstanceImageListRun

	// Server Type
	instanceServerTypeList := c.MustFind("instance", "server-type", "list")
	instanceServerTypeList.Run = customInstanceServerTypeListRun(instanceServerTypeList.Run)
}

type customSecurityGroupResponse struct {
	instance.SecurityGroup

	Rules []*instance.SecurityGroupRule
}

func customInstanceSecurityGroupGetRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	req := argsI.(*instance.GetSecurityGroupRequest)

	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)
	securityGroup, err := api.GetSecurityGroup(req)
	if err != nil {
		return nil, err
	}

	securityGroupRules, err := api.ListSecurityGroupRules(&instance.ListSecurityGroupRulesRequest{
		Zone:            req.Zone,
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

func customInstanceSecurityGroupDeleteRun(originalRun core.CommandRunner) core.CommandRunner {
	return func(ctx context.Context, argsI interface{}) (interface{}, error) {
		res, originalErr := originalRun(ctx, argsI)
		if originalErr == nil {
			return res, nil
		}

		if strings.HasSuffix(originalErr.Error(), "group is in use. you cannot delete it.") {
			req := argsI.(*instance.DeleteSecurityGroupRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			newError := &core.CliError{
				Err: fmt.Errorf("cannot delete security-group currently in use"),
			}

			// Get security-group.
			sg, err := api.GetSecurityGroup(&instance.GetSecurityGroupRequest{
				SecurityGroupID: req.SecurityGroupID,
			})
			if err != nil {
				// Ignore API error and return a minimal error.
				return nil, newError
			}

			// Create detail message.
			hint := "Attach all these instances to another security-group before deleting this one:"
			for _, s := range sg.SecurityGroup.Servers {
				hint += "\nscw instance server update server-id=" + s.ID + " security-group.id=$NEW_SECURITY_GROUP_ID"
			}

			newError.Hint = hint
			return nil, newError
		}

		return nil, originalErr
	}
}

// customInstanceImageListRun list the images for a given organization.
// A call to GetServer(..) with the ID contained in Image.FromServer retrieves more information about the server.
func customInstanceImageListRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
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

// customInstanceServerTypeListRun transforms the server map into a list to display a
// table of server types instead of a flat key/value list.
func customInstanceServerTypeListRun(originalRun core.CommandRunner) core.CommandRunner {
	return func(ctx context.Context, argsI interface{}) (interface{}, error) {

		type customServerType struct {
			Name            string        `json:"name"`
			MonthlyPrice    *scw.Money    `json:"monthly_price"`
			HourlyPrice     *scw.Money    `json:"hourly_price"`
			LocalVolumeSize scw.Size      `json:"local_volume_size"`
			CPU             uint32        `json:"cpu"`
			GPU             *uint64       `json:"gpu"`
			RAM             scw.Size      `json:"ram"`
			Arch            instance.Arch `json:"arch"`
		}

		originalRes, err := originalRun(ctx, argsI)
		if err != nil {
			return nil, err
		}

		listServersTypesResponse := originalRes.(*instance.ListServersTypesResponse)
		serverTypes := []*customServerType(nil)

		for name, serverType := range listServersTypesResponse.Servers {
			serverTypes = append(serverTypes, &customServerType{
				Name:            name,
				MonthlyPrice:    scw.NewMoneyFromFloat(float64(serverType.MonthlyPrice), "EUR", 2),
				HourlyPrice:     scw.NewMoneyFromFloat(float64(serverType.HourlyPrice), "EUR", 3),
				LocalVolumeSize: serverType.VolumesConstraint.MinSize,
				CPU:             serverType.Ncpus,
				GPU:             serverType.Gpu,
				RAM:             scw.Size(serverType.RAM),
				Arch:            serverType.Arch,
			})
		}

		sort.Slice(serverTypes, func(i, j int) bool {
			return serverTypes[i].Name < serverTypes[j].Name
		})

		return serverTypes, nil
	}
}
