package instance

import (
	"context"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	imageActionTimeout = 60 * time.Minute
)

//
// Marshalers
//

// imageStateMarshalSpecs allows to override the displayed instance.ImageState.
var (
	imageStateMarshalSpecs = human.EnumMarshalSpecs{
		instance.ImageStateCreating:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		instance.ImageStateAvailable: &human.EnumMarshalSpec{Attribute: color.FgGreen},
		instance.ImageStateError:     &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

func imagesMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type humanImage struct {
		ID               string
		Name             string
		State            instance.ImageState
		Public           bool
		Zone             scw.Zone
		Volumes          []scw.Size
		ServerName       string
		ServerID         string
		Arch             instance.Arch
		OrganizationID   string
		ProjectID        string
		CreationDate     *time.Time
		ModificationDate *time.Time
	}

	images := i.([]*imageListItem)
	humanImages := []*humanImage(nil)
	for _, image := range images {
		// For each image we want to display a list of volume size sepatated with `,`
		// e.g: 10 GB, 20 GB
		volumes := []scw.Size{
			image.RootVolume.Size,
		}
		// We must sort map key to make sure volume size are in the correct order.
		extraVolumeKeys := []string(nil)
		for key := range image.ExtraVolumes {
			extraVolumeKeys = append(extraVolumeKeys, key)
		}
		sort.Strings(extraVolumeKeys)

		for _, key := range extraVolumeKeys {
			volumes = append(volumes, image.ExtraVolumes[key].Size)
		}

		humanImages = append(humanImages, &humanImage{
			ID:               image.ID,
			Name:             image.Name,
			State:            image.State,
			Public:           image.Public,
			Zone:             image.Zone,
			Volumes:          volumes,
			ServerName:       image.ServerName,
			ServerID:         image.ServerID,
			Arch:             image.Arch,
			OrganizationID:   image.OrganizationID,
			ProjectID:        image.ProjectID,
			CreationDate:     image.CreationDate,
			ModificationDate: image.ModificationDate,
		})
	}
	return human.Marshal(humanImages, nil)
}

//
// Builders
//

// imageCreateBuilder overrides 'instance image create' to
// - rename extra-volumes arguments into additional-volumes
// - rename the argument 'root-volume' into 'snapshot-id'
func imageCreateBuilder(c *core.Command) *core.Command {
	type customCreateImageRequest struct {
		*instance.CreateImageRequest
		AdditionalVolumes []*instance.VolumeTemplate
		SnapshotID        string
		OrganizationID    *string
		ProjectID         *string
	}

	c.ArgSpecs.GetByName("extra-volumes.{key}.id").Short = "UUID of the snapshot to add"
	c.ArgSpecs.GetByName("extra-volumes.{key}.id").Name = "additional-snapshots.{index}.id"

	c.ArgSpecs.GetByName("extra-volumes.{key}.name").Short = "Name of the additional snapshot"
	c.ArgSpecs.GetByName("extra-volumes.{key}.name").Name = "additional-snapshots.{index}.name"

	c.ArgSpecs.GetByName("extra-volumes.{key}.size").Short = "Size of the additional snapshot"
	c.ArgSpecs.GetByName("extra-volumes.{key}.size").Name = "additional-snapshots.{index}.size"

	c.ArgSpecs.GetByName("extra-volumes.{key}.volume-type").Short = "Underlying volume type of the additional snapshot"
	c.ArgSpecs.GetByName("extra-volumes.{key}.volume-type").Name = "additional-snapshots.{index}.volume-type"

	c.ArgSpecs.GetByName("extra-volumes.{key}.organization").Short = "Organization ID that own the additional snapshot"
	c.ArgSpecs.GetByName("extra-volumes.{key}.organization").Name = "additional-snapshots.{index}.organization-id"

	c.ArgSpecs.GetByName("extra-volumes.{key}.project").Short = "Project ID that own the additional snapshot"
	c.ArgSpecs.GetByName("extra-volumes.{key}.project").Name = "additional-snapshots.{index}.project-id"

	c.ArgSpecs.GetByName("root-volume").Short = "UUID of the snapshot that will be used as root volume in the image"
	c.ArgSpecs.GetByName("root-volume").Name = "snapshot-id"

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreateImageRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customCreateImageRequest)

		request := args.CreateImageRequest
		request.RootVolume = args.SnapshotID
		request.ExtraVolumes = make(map[string]*instance.VolumeTemplate)
		request.Organization = args.OrganizationID
		request.Project = args.ProjectID

		// Extra volumes need to start at volumeIndex 1.
		volumeIndex := 1
		for _, volume := range args.AdditionalVolumes {
			request.ExtraVolumes[strconv.Itoa(volumeIndex)] = volume
			volumeIndex++
		}

		return runner(ctx, request)
	})

	return c
}

// customImage is based on instance.Image, with additional information about the server
type imageListItem struct {
	ID                string                      `json:"id"`
	Name              string                      `json:"name"`
	Arch              instance.Arch               `json:"arch"`
	CreationDate      *time.Time                  `json:"creation_date"`
	ModificationDate  *time.Time                  `json:"modification_date"`
	DefaultBootscript *instance.Bootscript        `json:"default_bootscript"`
	ExtraVolumes      map[string]*instance.Volume `json:"extra_volumes"`
	OrganizationID    string                      `json:"organization"`
	ProjectID         string                      `json:"project"`
	Public            bool                        `json:"public"`
	RootVolume        *instance.VolumeSummary     `json:"root_volume"`
	State             instance.ImageState         `json:"state"`

	// Replace Image.FromServer
	ServerID   string   `json:"server_id"`
	ServerName string   `json:"server_name"`
	Zone       scw.Zone `json:"zone"`
}

// imageListBuilder list the images for a given organization/project.
// A call to GetServer(..) with the ID contained in Image.FromServer retrieves more information about the server.
func imageListBuilder(c *core.Command) *core.Command {
	type customListImageRequest struct {
		*instance.ListImagesRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListImageRequest{})

	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		// Get images
		args := argsI.(*customListImageRequest)

		if args.ListImagesRequest == nil {
			args.ListImagesRequest = &instance.ListImagesRequest{}
		}

		req := args.ListImagesRequest
		req.Organization = args.OrganizationID
		req.Project = args.ProjectID
		req.Public = scw.BoolPtr(false)
		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		listImagesResponse, err := api.ListImages(req, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		images := listImagesResponse.Images

		// Builds customImages
		customImages := []*imageListItem(nil)
		for _, image := range images {
			newCustomImage := &imageListItem{
				ID:                image.ID,
				Name:              image.Name,
				Arch:              image.Arch,
				CreationDate:      image.CreationDate,
				ModificationDate:  image.ModificationDate,
				DefaultBootscript: image.DefaultBootscript,
				ExtraVolumes:      image.ExtraVolumes,
				OrganizationID:    image.Organization,
				ProjectID:         image.Project,
				Public:            image.Public,
				RootVolume:        image.RootVolume,
				State:             image.State,
				Zone:              image.Zone,
			}
			customImages = append(customImages, newCustomImage)

			if image.FromServer == "" {
				continue
			}

			serverReq := instance.GetServerRequest{
				Zone:     req.Zone,
				ServerID: image.FromServer,
			}
			getServerResponse, err := api.GetServer(&serverReq)
			if _, ok := err.(*scw.ResourceNotFoundError); ok {
				newCustomImage.ServerName = "-"
				continue
			}
			if err != nil {
				return nil, err
			}
			newCustomImage.ServerID = getServerResponse.Server.ID
			newCustomImage.ServerName = getServerResponse.Server.Name
		}

		return customImages, nil
	}

	return c
}

// imageDeleteBuilder override delete command to:
//  - add a with-snapshots parameter
func imageDeleteBuilder(c *core.Command) *core.Command {
	type customDeleteImageRequest struct {
		*instance.DeleteImageRequest
		WithSnapshots bool
	}

	c.ArgsType = reflect.TypeOf(customDeleteImageRequest{})
	c.ArgSpecs.AddBefore("zone", &core.ArgSpec{
		Name:  "with-snapshots",
		Short: "Delete the snapshots attached to this image",
	})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customDeleteImageRequest)

		api := instance.NewAPI(core.ExtractClient(ctx))

		// If we want to delete snapshot we must GET image before we delete it
		image := (*instance.Image)(nil)
		if args.WithSnapshots {
			res, err := api.GetImage(&instance.GetImageRequest{
				Zone:    args.Zone,
				ImageID: args.ImageID,
			})
			if err != nil {
				return nil, err
			}
			image = res.Image
		}

		// Call the generated delete
		runnerRes, err := runner(ctx, args.DeleteImageRequest)
		if err != nil {
			return nil, err
		}

		// Once the image is deleted we can delete snapshots.
		if args.WithSnapshots {
			snapshotIDs := []string{
				image.RootVolume.ID,
			}
			for _, snapshot := range image.ExtraVolumes {
				snapshotIDs = append(snapshotIDs, snapshot.ID)
			}
			for _, snapshotID := range snapshotIDs {
				err := api.DeleteSnapshot(&instance.DeleteSnapshotRequest{
					Zone:       args.Zone,
					SnapshotID: snapshotID,
				})
				if err != nil {
					return nil, err
				}
			}
		}
		return runnerRes, nil
	})

	return c
}

func imageWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for image to reach a stable state`,
		Long:      `Wait for image to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the image.`,
		Namespace: "instance",
		Resource:  "image",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(instance.WaitForImageRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := instance.NewAPI(core.ExtractClient(ctx))
			return api.WaitForImage(&instance.WaitForImageRequest{
				Zone:          argsI.(*instance.WaitForImageRequest).Zone,
				ImageID:       argsI.(*instance.WaitForImageRequest).ImageID,
				Timeout:       scw.TimeDurationPtr(imageActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `ID of the image.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a image to reach a stable state",
				ArgsJSON: `{"image_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
