package instance

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

type instanceCreateServerRequest struct {
	Zone              scw.Zone
	ProjectID         *string
	Image             string
	Type              string
	Name              string
	RootVolume        string
	AdditionalVolumes []string
	IP                string
	Tags              []string
	IPv6              bool
	Stopped           bool
	SecurityGroupID   string
	PlacementGroupID  string
	BootscriptID      string
	CloudInit         string
	BootType          string

	// Deprecated, use project-id instead
	OrganizationID *string
}

func serverCreateCommand() *core.Command {
	return &core.Command{
		Short:     `Create server`,
		Long:      `Create an instance server.`,
		Namespace: "instance",
		Verb:      "create",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(instanceCreateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "image",
				Short:            "Image ID or label of the server",
				Default:          core.DefaultValueSetter("ubuntu_focal"),
				Required:         true,
				AutoCompleteFunc: instanceServerCreateImageAutoCompleteFunc,
			},
			{
				Name:    "type",
				Short:   "Server commercial type",
				Default: core.DefaultValueSetter("DEV1-S"),
				EnumValues: []string{
					"GP1-XS",
					"GP1-S",
					"GP1-M",
					"GP1-L",
					"GP1-XL",
					"DEV1-S",
					"DEV1-M",
					"DEV1-L",
					"DEV1-XL",
					"RENDER-S",
					"STARDUST1-S",
					"ENT1-S",
					"ENT1-M",
					"ENT1-L",
					"ENT1-XL",
					"ENT1-2XL",
				},
				ValidateFunc: func(argSpec *core.ArgSpec, value interface{}) error {
					// Allow all commercial types
					return nil
				},
			},
			{
				Name:    "name",
				Short:   "Server name",
				Default: core.RandomValueGenerator("srv"),
			},
			{
				Name:  "root-volume",
				Short: "Local root volume of the server",
			},
			{
				Name:  "additional-volumes.{index}",
				Short: "Additional local and block volumes attached to your server",
			},
			{
				Name:    "ip",
				Short:   `Either an IP, an IP ID, 'new' to create a new IP, 'dynamic' to use a dynamic IP or 'none' for no public IP (new | dynamic | none | <id> | <address>)`,
				Default: core.DefaultValueSetter("new"),
			},
			{
				Name:  "tags.{index}",
				Short: "Server tags",
			},
			{
				Name:  "ipv6",
				Short: "Enable IPv6",
			},
			{
				Name:  "stopped",
				Short: "Do not start server after its creation",
			},
			{
				Name:  "security-group-id",
				Short: "The security group ID it use for this server",
			},
			{
				Name:  "placement-group-id",
				Short: "The placement group ID in witch the server has to be created",
			},
			{
				Name:  "bootscript-id",
				Short: "The bootscript ID to use, if empty the local boot will be used",
			},
			{
				Name:        "cloud-init",
				Short:       "The cloud-init script to use",
				CanLoadFile: true,
			},
			{
				Name:       "boot-type",
				Short:      "The boot type to use, if empty the local boot will be used. Will be overwritten to bootscript if bootscript-id is set.",
				Default:    core.DefaultValueSetter(instance.BootTypeLocal.String()),
				EnumValues: []string{instance.BootTypeLocal.String(), instance.BootTypeBootscript.String(), instance.BootTypeRescue.String()},
			},
			core.ProjectIDArgSpec(),
			core.ZoneArgSpec(),
			core.OrganizationIDArgSpec(),
		},
		Run:      instanceServerCreateRun,
		WaitFunc: instanceWaitServerCreateRun(),
		SeeAlsos: []*core.SeeAlso{{
			Short:   "List marketplace label images",
			Command: "scw marketplace image list",
		}},
		Examples: []*core.Example{
			{
				Short:    "Create and start an instance on Ubuntu Focal",
				ArgsJSON: `{"image":"ubuntu_focal","start":true}`,
			},
			{
				Short:    "Create a GP1-XS instance, give it a name and add tags",
				ArgsJSON: `{"image":"ubuntu_focal","type":"GP1-XS","name":"foo","tags":["prod","blue"]}`,
			},
			{
				Short:    "Create an instance with 2 additional block volumes (50GB and 100GB)",
				ArgsJSON: `{"image":"ubuntu_focal","additional_volumes":["block:50GB","block:100GB"]}`,
			},
			{
				Short:    "Create an instance with 2 local volumes (10GB and 10GB)",
				ArgsJSON: `{"image":"ubuntu_focal","root_volume":"local:10GB","additional_volumes":["local:10GB"]}`,
			},
			{
				Short: "Use an existing IP",
				Raw: `ip=$(scw instance ip create | grep id | awk '{ print $2 }')
scw instance server create image=ubuntu_focal ip=$ip`,
			},
		},
	}
}

func instanceWaitServerCreateRun() core.WaitFunc {
	return func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		return instance.NewAPI(core.ExtractClient(ctx)).WaitForServer(&instance.WaitForServerRequest{
			Zone:          argsI.(*instanceCreateServerRequest).Zone,
			ServerID:      respI.(*instance.Server).ID,
			Timeout:       scw.TimeDurationPtr(serverActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}
}

func instanceServerCreateRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*instanceCreateServerRequest)

	//
	// STEP 1: Argument validation and API requests creation.
	//

	needIPCreation := false

	serverReq := &instance.CreateServerRequest{
		Zone:           args.Zone,
		Organization:   args.OrganizationID,
		Project:        args.ProjectID,
		Name:           args.Name,
		CommercialType: args.Type,
		EnableIPv6:     args.IPv6,
		Tags:           args.Tags,
	}

	client := core.ExtractClient(ctx)
	apiMarketplace := marketplace.NewAPI(client)
	apiInstance := instance.NewAPI(client)

	//
	// Image.
	//
	// Could be:
	// - A local image UUID
	// - An image label
	//
	switch {
	case !validation.IsUUID(args.Image):
		// Find the corresponding local image UUID.
		imageID, err := apiMarketplace.GetLocalImageIDByLabel(&marketplace.GetLocalImageIDByLabelRequest{
			Zone:           args.Zone,
			ImageLabel:     args.Image,
			CommercialType: serverReq.CommercialType,
		})
		if err != nil {
			return nil, fmt.Errorf("bad image label '%s' for %s", args.Image, serverReq.CommercialType)
		}
		serverReq.Image = imageID
	default:
		serverReq.Image = args.Image
	}

	getImageResponse, err := apiInstance.GetImage(&instance.GetImageRequest{
		Zone:    args.Zone,
		ImageID: serverReq.Image,
	})
	if err != nil {
		logger.Warningf("cannot get image %s: %s", serverReq.Image, err)
	}

	serverType := getServeType(apiInstance, serverReq.Zone, serverReq.CommercialType)

	if serverType != nil && getImageResponse != nil {
		if err := validateImageServerTypeCompatibility(getImageResponse.Image, serverType, serverReq.CommercialType); err != nil {
			return nil, err
		}
	} else {
		logger.Warningf("skipping image server-type compatibility validation")
	}

	//
	// IP.
	//
	// Could be:
	// - "new"
	// - A flexible IP UUID
	// - A flexible IP address
	// - "dynamic"
	// - "none"
	//
	switch {
	case args.IP == "", args.IP == "new":
		needIPCreation = true
	case validation.IsUUID(args.IP):
		serverReq.PublicIP = scw.StringPtr(args.IP)
	case net.ParseIP(args.IP) != nil:
		// Find the corresponding flexible IP UUID.
		logger.Debugf("finding public IP UUID from address: %s", args.IP)
		res, err := apiInstance.GetIP(&instance.GetIPRequest{
			Zone: args.Zone,
			IP:   args.IP,
		})
		if err != nil { // FIXME: isNotFoundError
			return nil, fmt.Errorf("%s does not belong to you", args.IP)
		}
		serverReq.PublicIP = scw.StringPtr(res.IP.ID)
	case args.IP == "dynamic":
		serverReq.DynamicIPRequired = scw.BoolPtr(true)
	case args.IP == "none":
		serverReq.DynamicIPRequired = scw.BoolPtr(false)
	default:
		return nil, fmt.Errorf(`invalid IP "%s", should be either 'new', 'dynamic', 'none', an IP address ID or a reserved flexible IP address`, args.IP)
	}

	//
	// Volumes.
	//
	// More format details in buildVolumeTemplate function.
	//
	if len(args.AdditionalVolumes) > 0 || args.RootVolume != "" {
		// Create initial volume template map.
		volumes, err := buildVolumes(apiInstance, args.Zone, serverReq.Name, args.RootVolume, args.AdditionalVolumes)
		if err != nil {
			return nil, err
		}

		// Validate root volume type and size.
		if getImageResponse != nil {
			if err := validateRootVolume(getImageResponse.Image.RootVolume.Size, volumes["0"]); err != nil {
				return nil, err
			}
		} else {
			logger.Warningf("skipping root volume validation")
		}

		// Validate total local volume sizes.
		if serverType != nil {
			if err := validateLocalVolumeSizes(volumes, serverType, serverReq.CommercialType); err != nil {
				return nil, err
			}
		} else {
			logger.Warningf("skip local volume size validation")
		}

		// Sanitize the volume map to respect API schemas
		serverReq.Volumes = sanitizeVolumeMap(serverReq.Name, volumes)
	}

	//
	// BootType.
	//
	bootType := instance.BootType(args.BootType)
	serverReq.BootType = &bootType

	//
	// Bootscript.
	//
	if args.BootscriptID != "" {
		if !validation.IsUUID(args.BootscriptID) {
			return nil, fmt.Errorf("bootscript ID %s is not a valid UUID", args.BootscriptID)
		}
		_, err := apiInstance.GetBootscript(&instance.GetBootscriptRequest{
			Zone:         args.Zone,
			BootscriptID: args.BootscriptID,
		})
		if err != nil { // FIXME: isNotFoundError
			return nil, fmt.Errorf("bootscript ID %s does not exist", args.BootscriptID)
		}

		serverReq.Bootscript = scw.StringPtr(args.BootscriptID)
		bootType := instance.BootTypeBootscript
		serverReq.BootType = &bootType
	}

	//
	// Security Group.
	//
	if args.SecurityGroupID != "" {
		serverReq.SecurityGroup = scw.StringPtr(args.SecurityGroupID)
	}

	//
	// Placement Group.
	//
	if args.PlacementGroupID != "" {
		serverReq.PlacementGroup = scw.StringPtr(args.PlacementGroupID)
	}

	//
	// STEP 2: Resource creations and modifications.
	//

	//
	// IP
	//
	if needIPCreation {
		logger.Debugf("creating IP")

		res, err := apiInstance.CreateIP(&instance.CreateIPRequest{
			Zone:         args.Zone,
			Project:      args.ProjectID,
			Organization: args.OrganizationID,
		})
		if err != nil {
			return nil, fmt.Errorf("error while creating your public IP: %s", err)
		}
		serverReq.PublicIP = scw.StringPtr(res.IP.ID)
		logger.Debugf("IP created: %s", serverReq.PublicIP)
	}

	//
	// Server Creation
	//
	logger.Debugf("creating server")
	serverRes, err := apiInstance.CreateServer(serverReq)
	if err != nil {
		if needIPCreation && serverReq.PublicIP != nil {
			// Delete the created IP
			logger.Debugf("deleting created IP: %s", serverReq.PublicIP)
			err := apiInstance.DeleteIP(&instance.DeleteIPRequest{
				Zone: args.Zone,
				IP:   *serverReq.PublicIP,
			})
			if err != nil {
				logger.Warningf("cannot delete the create IP %s: %s.", serverReq.PublicIP, err)
			}
		}

		return nil, fmt.Errorf("cannot create the server: %s", err)
	}
	server := serverRes.Server
	logger.Debugf("server created %s", server.ID)

	//
	// Cloud-init
	//
	if args.CloudInit != "" {
		err := apiInstance.SetServerUserData(&instance.SetServerUserDataRequest{
			Zone:     args.Zone,
			ServerID: server.ID,
			Key:      "cloud-init",
			Content:  bytes.NewBufferString(args.CloudInit),
		})
		if err != nil {
			logger.Warningf("error while setting up your cloud-init metadata: %s. Note that the server is successfully created.", err)
		} else {
			logger.Debugf("cloud-init set")
		}
	}

	//
	// Start server by default
	//
	if !args.Stopped {
		logger.Debugf("starting server")
		_, err := apiInstance.ServerAction(&instance.ServerActionRequest{
			Zone:     args.Zone,
			ServerID: server.ID,
			Action:   instance.ServerActionPoweron,
		})
		if err != nil {
			logger.Warningf("Cannot start the server: %s. Note that the server is successfully created.", err)
		} else {
			logger.Debugf("server started")
		}
	}

	return server, nil
}

// buildVolumes creates the initial volume map.
// It is not the definitive one, it will be mutated all along the process.
func buildVolumes(api *instance.API, zone scw.Zone, serverName, rootVolume string, additionalVolumes []string) (map[string]*instance.VolumeServerTemplate, error) {
	volumes := make(map[string]*instance.VolumeServerTemplate)
	if rootVolume != "" {
		rootVolumeTemplate, err := buildVolumeTemplate(api, zone, rootVolume)
		if err != nil {
			return nil, err
		}

		volumes["0"] = rootVolumeTemplate
	}

	for i, v := range additionalVolumes {
		volumeTemplate, err := buildVolumeTemplate(api, zone, v)
		if err != nil {
			return nil, err
		}
		index := strconv.Itoa(i + 1)
		volumeTemplate.Name = serverName + "-" + index

		// Remove extra data for API validation.
		if volumeTemplate.ID != "" {
			volumeTemplate = &instance.VolumeServerTemplate{
				ID:   volumeTemplate.ID,
				Name: volumeTemplate.Name,
			}
		}

		volumes[index] = volumeTemplate
	}

	return volumes, nil
}

// buildVolumeTemplate creates a instance.VolumeTemplate from a 'volumes' argument item.
//
// Volumes definition must be through multiple arguments (eg: volumes.0="l:20GB" volumes.1="b:100GB")
//
// A valid volume format is either
// - a "creation" format: ^((local|l|block|b):)?\d+GB?$ (size is handled by go-humanize, so other sizes are supported)
// - a UUID format
//
func buildVolumeTemplate(api *instance.API, zone scw.Zone, flagV string) (*instance.VolumeServerTemplate, error) {
	parts := strings.Split(strings.TrimSpace(flagV), ":")

	// Create volume.
	if len(parts) == 2 {
		vt := &instance.VolumeServerTemplate{}

		switch parts[0] {
		case "l", "local":
			vt.VolumeType = instance.VolumeVolumeTypeLSSD
		case "b", "block":
			vt.VolumeType = instance.VolumeVolumeTypeBSSD
		default:
			return nil, fmt.Errorf("invalid volume type %s in %s volume", parts[0], flagV)
		}

		size, err := humanize.ParseBytes(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid size format %s in %s volume", parts[1], flagV)
		}
		vt.Size = scw.Size(size)

		return vt, nil
	}

	// UUID format.
	if len(parts) == 1 && validation.IsUUID(parts[0]) {
		return buildVolumeTemplateFromUUID(api, zone, parts[0])
	}

	return nil, &core.CliError{
		Err:     fmt.Errorf("invalid volume format '%s'", flagV),
		Details: "",
		Hint:    `You must provide either a UUID ("11111111-1111-1111-1111-111111111111"), a local volume size ("local:100G" or "l:100G") or a block volume size ("block:100G" or "b:100G").`,
	}
}

// buildVolumeTemplateFromUUID validate an UUID volume and add their types and sizes.
// Add volume types and sizes allow US to treat UUID volumes like the others and simplify the implementation.
// The instance API refuse the type and the size for UUID volumes, therefore,
// buildVolumeMap function will remove them.
func buildVolumeTemplateFromUUID(api *instance.API, zone scw.Zone, volumeUUID string) (*instance.VolumeServerTemplate, error) {
	res, err := api.GetVolume(&instance.GetVolumeRequest{
		Zone:     zone,
		VolumeID: volumeUUID,
	})
	if err != nil { // FIXME: isNotFoundError
		return nil, fmt.Errorf("volume %s does not exist", volumeUUID)
	}

	// Check that volume is not already attached to a server.
	if res.Volume.Server != nil {
		return nil, fmt.Errorf("volume %s is already attached to %s server", res.Volume.ID, res.Volume.Server.ID)
	}

	return &instance.VolumeServerTemplate{
		ID:         res.Volume.ID,
		VolumeType: res.Volume.VolumeType,
		Size:       res.Volume.Size,
	}, nil
}

func validateImageServerTypeCompatibility(image *instance.Image, serverType *instance.ServerType, CommercialType string) error {
	// An instance might not have any constraints on the local volume size
	if serverType.VolumesConstraint.MaxSize == 0 {
		return nil
	}
	if image.RootVolume.Size > serverType.VolumesConstraint.MaxSize {
		return fmt.Errorf("image %s requires %s on root volume, but root volume is constrained between %s and %s on %s",
			image.ID,
			humanize.Bytes(uint64(image.RootVolume.Size)),
			humanize.Bytes(uint64(serverType.VolumesConstraint.MinSize)),
			humanize.Bytes(uint64(serverType.VolumesConstraint.MaxSize)),
			CommercialType,
		)
	}

	return nil
}

// validateLocalVolumeSizes validates the total size of local volumes.
func validateLocalVolumeSizes(volumes map[string]*instance.VolumeServerTemplate, serverType *instance.ServerType, commercialType string) error {
	// Calculate local volume total size.
	var localVolumeTotalSize scw.Size
	for _, volume := range volumes {
		if volume.VolumeType == instance.VolumeVolumeTypeLSSD {
			localVolumeTotalSize += volume.Size
		}
	}

	volumeConstraint := serverType.VolumesConstraint

	// If no root volume provided, count the default root volume size added by the API.
	if rootVolume := volumes["0"]; rootVolume == nil {
		localVolumeTotalSize += volumeConstraint.MinSize
	}

	if localVolumeTotalSize < volumeConstraint.MinSize || localVolumeTotalSize > volumeConstraint.MaxSize {
		min := humanize.Bytes(uint64(volumeConstraint.MinSize))
		if volumeConstraint.MinSize == volumeConstraint.MaxSize {
			return fmt.Errorf("%s total local volume size must be equal to %s", commercialType, min)
		}

		max := humanize.Bytes(uint64(volumeConstraint.MaxSize))
		return fmt.Errorf("%s total local volume size must be between %s and %s", commercialType, min, max)
	}

	return nil
}

func validateRootVolume(imageRequiredSize scw.Size, rootVolume *instance.VolumeServerTemplate) error {
	if rootVolume == nil {
		return nil
	}

	if rootVolume.ID != "" {
		return &core.CliError{
			Err:     fmt.Errorf("you cannot use an existing volume as a root volume"),
			Details: "You must create an image of this volume and use its ID in the 'image' argument.",
		}
	}

	if rootVolume.Size < imageRequiredSize {
		return fmt.Errorf("first volume size must be at least %s for this image", humanize.Bytes(uint64(imageRequiredSize)))
	}

	return nil
}

// sanitizeVolumeMap removes extra data for API validation.
func sanitizeVolumeMap(serverName string, volumes map[string]*instance.VolumeServerTemplate) map[string]*instance.VolumeServerTemplate {
	m := make(map[string]*instance.VolumeServerTemplate)

	for index, v := range volumes {
		v.Name = serverName + "-" + index

		// Remove extra data for API validation.
		switch {
		case v.ID != "":
			v = &instance.VolumeServerTemplate{
				ID:   v.ID,
				Name: v.Name,
			}
		case index == "0" && v.Size != 0:
			v = &instance.VolumeServerTemplate{
				VolumeType: v.VolumeType,
				Size:       v.Size,
			}
		}
		m[index] = v
	}

	return m
}

func instanceServerCreateImageAutoCompleteFunc(ctx context.Context, prefix string) core.AutocompleteSuggestions {
	suggestions := core.AutocompleteSuggestions(nil)

	client := core.ExtractClient(ctx)
	api := marketplace.NewAPI(client)

	res, err := api.ListImages(&marketplace.ListImagesRequest{}, scw.WithAllPages())
	if err != nil {
		return nil
	}

	prefix = strings.ToLower(strings.Replace(prefix, "-", "_", -1))

	for _, image := range res.Images {
		if strings.HasPrefix(image.Label, prefix) {
			suggestions = append(suggestions, image.Label)
		}
	}

	return suggestions
}

// getServeType is a util to get a instance.ServerType by its commercialType
func getServeType(apiInstance *instance.API, zone scw.Zone, commercialType string) *instance.ServerType {
	serverType := (*instance.ServerType)(nil)

	serverTypesRes, err := apiInstance.ListServersTypes(&instance.ListServersTypesRequest{
		Zone: zone,
	})
	if err != nil {
		logger.Warningf("cannot get server types: %s", err)
	} else {
		serverType = serverTypesRes.Servers[commercialType]
		if serverType == nil {
			logger.Warningf("unrecognized server type: %s", commercialType)
		}
	}

	return serverType
}
