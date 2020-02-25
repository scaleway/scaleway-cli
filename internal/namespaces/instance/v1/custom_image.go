package instance

import (
	"context"
	"reflect"
	"strconv"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Builders
//

// imageCreateBuilder overrides 'instance image create' to
// - rename extra-volumes arguments into additional-volumes
// - rename the argument 'root-volume' into 'snapshot-id'
func imageCreateBuilder(c *core.Command) *core.Command {

	type customCreateImageRequest struct {
		*instance.CreateImageRequest
		AdditionalVolumes map[string]*instance.VolumeTemplate
		SnapshotID        string
	}

	c.ArgSpecs.GetByName("extra-volumes.{key}.id").Name = "additional-volumes.{key}.id"
	c.ArgSpecs.GetByName("extra-volumes.{key}.name").Name = "additional-volumes.{key}.name"
	c.ArgSpecs.GetByName("extra-volumes.{key}.size").Name = "additional-volumes.{key}.size"
	c.ArgSpecs.GetByName("extra-volumes.{key}.volume-type").Name = "additional-volumes.{key}.volume-type"
	c.ArgSpecs.GetByName("extra-volumes.{key}.organization").Name = "additional-volumes.{key}.organization"

	c.ArgSpecs.GetByName("root-volume").Name = "snapshot-id"

	c.ArgsType = reflect.TypeOf(customCreateImageRequest{})

	originalRun := c.Run

	c.Run = func(ctx context.Context, args interface{}) (i interface{}, e error) {
		args := argsI.(*customCreateImageRequest)

		request := args.CreateImageRequest
		request.RootVolume = args.SnapshotID
		request.ExtraVolumes = make(map[string]*instance.VolumeTemplate)

		// Extra volumes need to start at volumeIndex 1.
		volumeIndex := 1
		for _, volume := range args.AdditionalVolumes {
			request.ExtraVolumes[strconv.Itoa(volumeIndex)] = volume
			volumeIndex++
		}

		return originalRun(ctx, request)
	}

	return c
}

// imageListBuilder list the images for a given organization.
// A call to GetServer(..) with the ID contained in Image.FromServer retrieves more information about the server.
func imageListBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
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
			newCustomImage := &customImage{
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
				Zone:              image.Zone,
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
				newCustomImage.ServerID = getServerResponse.Server.ID
				newCustomImage.ServerName = getServerResponse.Server.Name
			}
			customImages = append(customImages, newCustomImage)
		}

		return customImages, nil
	}

	return c
}
