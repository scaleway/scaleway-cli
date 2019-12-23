package instance

import (
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

// TODO: Add cloud-init
type instanceCreateServerRequest struct {
	Image            string
	Type             string
	Name             string
	Volumes          []string
	IP               string
	Tags             []string
	IPv6             bool
	Start            bool
	SecurityGroupID  string
	PlacementGroupID string
	BootscriptID     string
}

// TODO: Remove all error uppercase and punctuations when [APIGW-1367] will be done
func instanceServerCreate() *core.Command {
	return &core.Command{
		Short:     `Create server`,
		Long:      `Create an instance server.`, // TODO: Add examples [APIGW-1371]
		Namespace: "instance",
		Verb:      "create",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(instanceCreateServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "image",
				Short:            "Image ID or label of the server",
				Required:         true,
				AutoCompleteFunc: instanceServerCreateImageAutoCompleteFunc,
			},
			{
				Name:       "type",
				Short:      "Server commercial type",
				Default:    core.DefaultValueSetter("DEV1-S"),
				EnumValues: []string{"GP1-XS", "GP1-S", "GP1-M", "GP1-L", "GP1-XL", "DEV1-S", "DEV1-M", "DEV1-L", "DEV1-XL", "RENDER-S"},
			},
			{
				Name:    "name",
				Short:   "Server name",
				Default: core.RandomValueGenerator("srv"),
			},
			{
				Name:  "volumes.{idx}",
				Short: "Local and short volumes attached to your server", // TODO: Add examples [APIGW-1371]
			},
			{
				Name:       "ip",
				Short:      "Either an IP, an IP ID, 'new' to create a new IP, 'dynamic' to use a dynamic IP or 'none' for no public IP",
				Default:    core.DefaultValueSetter("new"),
				EnumValues: []string{"new", "dynamic", "none", "<id>", "<address>"},
			},
			{
				Name:  "tags.{idx}",
				Short: "Server tags",
			},
			{
				Name:  "ipv6",
				Short: "Enable IPv6",
			},
			{
				Name:  "start",
				Short: "Start the server after its creation",
			},
			{
				Name:  "security-group-id",
				Short: "The security group ID it use for this server",
			},
			{
				Name:  "placement-group-id",
				Short: "The security group ID in witch the server has to be created",
			},
			{
				Name:  "bootscript-id",
				Short: "The bootscript ID to use, if empty the local boot will be used",
			},
		},
		Run:      instanceServerCreateRun,
		WaitFunc: instanceWaitServerCreateRun,
		SeeAlsos: []*core.SeeAlso{{
			Short:   "List marketplace label images",
			Command: "scw marketplace image list",
		}},
		Examples: []*core.Example{}, // TODO: Add examples [APIGW-1371]
	}
}

func instanceWaitServerCreateRun(ctx context.Context, _, respI interface{}) error {
	_, err := instance.NewAPI(core.ExtractClient(ctx)).WaitForServer(&instance.WaitForServerRequest{
		ServerID: respI.(*instance.Server).ID,
		Timeout:  serverActionTimeout,
	})
	return err
}

func instanceServerCreateRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	args := argsI.(*instanceCreateServerRequest)

	//
	// STEP 1: Argument validation and API requests creation.
	//

	needIPCreation := false

	serverReq := &instance.CreateServerRequest{
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
		imageId, err := apiMarketplace.GetLocalImageIDByLabel(&marketplace.GetLocalImageIDByLabelRequest{
			ImageLabel:     args.Image,
			CommercialType: serverReq.CommercialType,
		})
		if err != nil {
			return nil, fmt.Errorf("Bad image label '%s' for %s.", args.Image, serverReq.CommercialType)
		}
		serverReq.Image = imageId
	default:
		serverReq.Image = args.Image
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
		logger.Infof("finding public IP UUID from address: %s", args.IP)
		res, err := apiInstance.GetIP(&instance.GetIPRequest{IP: args.IP})
		if err != nil { // FIXME: isNotFoundError
			return nil, fmt.Errorf("%s does not belongs to you.", args.IP)
		}
		serverReq.PublicIP = scw.StringPtr(res.IP.ID)
	case args.IP == "dynamic":
		serverReq.DynamicIPRequired = scw.BoolPtr(true)
	case args.IP == "none":
	default:
		return nil, fmt.Errorf(`Invalid IP "%s", should be either 'new', 'dynamic', 'none', an IP address ID or a reserved flexible IP address.`, args.IP)
	}

	//
	// Volumes.
	//
	// More format details in buildVolumeTemplate function.
	//
	if len(args.Volumes) > 0 {
		// Get organization ID.
		organizationID := core.GetOrganizationIdFromContext(ctx)

		// Create initial volume template map.
		volumes, err := createInitialVolumeMap(apiInstance, organizationID, args.Volumes)
		if err != nil {
			return nil, err
		}

		// Validate root volume type and size.
		if err := validateRootVolume(apiInstance, serverReq.Image, volumes[0]); err != nil {
			return nil, err
		}

		// Validate total local volume sizes.
		if err := validateLocalVolumeSizes(apiInstance, volumes, serverReq.CommercialType); err != nil {
			return nil, err
		}

		// Create final volume map.
		serverReq.Volumes = buildVolumeMap(serverReq.Name, volumes)
	}

	//
	// Bootscript.
	//
	if args.BootscriptID != "" {
		if !validation.IsUUID(args.BootscriptID) {
			return nil, fmt.Errorf("Bootscript ID %s is not a valid UUID.", args.BootscriptID)
		}
		_, err := apiInstance.GetBootscript(&instance.GetBootscriptRequest{BootscriptID: args.BootscriptID})
		if err != nil { // FIXME: isNotFoundError
			return nil, fmt.Errorf("Bootscript ID %s does not exists.", args.BootscriptID)
		}

		serverReq.Bootscript = scw.StringPtr(args.BootscriptID)
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
		logger.Infof("creating IP")
		res, err := apiInstance.CreateIP(&instance.CreateIPRequest{})
		if err != nil {
			return nil, fmt.Errorf("Error while creating your public IP: %s.", err)
		}
		serverReq.PublicIP = scw.StringPtr(res.IP.ID)
		logger.Infof("IP created: %s", serverReq.PublicIP)
	}

	//
	// Server
	//
	logger.Infof("creating server")
	serverRes, err := apiInstance.CreateServer(serverReq)
	if err != nil {
		if needIPCreation && serverReq.PublicIP != nil {
			// Delete the created IP
			logger.Infof("deleting created IP: %s", serverReq.PublicIP)
			err := apiInstance.DeleteIP(&instance.DeleteIPRequest{IP: *serverReq.PublicIP})
			if err != nil {
				logger.Warningf("cannot delete the create IP %s: %s.", serverReq.PublicIP, err)
			}
		}

		return nil, fmt.Errorf("Cannot create the server: %s.", err)
	}
	server := serverRes.Server
	logger.Infof("server created %s", server.ID)

	//
	// Start
	//
	if args.Start {
		// TODO: Use the wait flag when it will be implemented [APIGW-1313]
		logger.Infof("starting server")
		_, err := apiInstance.ServerAction(&instance.ServerActionRequest{
			ServerID: server.ID,
			Action:   instance.ServerActionPoweron,
		})
		if err != nil {
			logger.Warningf("Cannot start the server: %s. Note that the server is successfully created.", err)
		} else {
			logger.Infof("server started")
		}
	}

	return server, nil
}

// createInitialVolumeMap creates the initial volume map.
// It is not the definitive one, it will be mutated all along the process.
func createInitialVolumeMap(api *instance.API, organizationID string, argsVolumes []string) ([]*instance.VolumeTemplate, error) {
	var volumes []*instance.VolumeTemplate

	for _, v := range argsVolumes {
		for _, flagV := range strings.Split(v, ",") {
			flagV = strings.TrimSpace(flagV)
			vt, err := buildVolumeTemplate(api, organizationID, flagV)
			if err != nil {
				return nil, err
			}

			volumes = append(volumes, vt)
		}
	}

	return volumes, nil
}

// buildVolumeTemplate creates a instance.VolumeTemplate from a 'volumes' argument item.
//
// Volumes definition could be:
// - multiple arguments (eg: volumes="l:20GB" volumes="b:100GB")
// - a single argument representing an array (eg: volumes="l:20GB, b:100GB")
// - a mix (eg: volumes="l:10GB, local:10GB" volumes="b:100GB")
//
// A valid volume format is either
// - a "creation" format: ^((local|l|block|b):)?\d+GB?$ (size is handled by go-humanize, so other sizes are supported)
// - a UUID format
//
func buildVolumeTemplate(api *instance.API, orgID, flagV string) (*instance.VolumeTemplate, error) {
	parts := strings.Split(flagV, ":")

	// Create volume.
	if len(parts) == 2 {
		vt := &instance.VolumeTemplate{}

		switch parts[0] {
		case "l", "local":
			vt.VolumeType = instance.VolumeTypeLSSD
		case "b", "block":
			vt.VolumeType = instance.VolumeTypeBSSD
		default:
			return nil, fmt.Errorf("Invalid volume type %s in %s volume.", parts[0], flagV)
		}

		size, err := humanize.ParseBytes(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Invalid size format %s in %s volume.", parts[1], flagV) // TODO: improve msg [APIGW-1371]
		}
		vt.Size = scw.Size(size)

		vt.Organization = orgID

		return vt, nil
	}

	// UUID format.
	if len(parts) == 1 && validation.IsUUID(parts[0]) {
		return buildVolumeTemplateFromUUID(api, parts[0])
	}

	return nil, fmt.Errorf("Invalid volume format '%s'.", flagV) // TODO: improve msg [APIGW-1371]
}

// buildVolumeTemplateFromUUID validate an UUID volume and add their types and sizes.
// Add volume types and sizes allow US to treat UUID volumes like the others and simplify the implementation.
// The instance API refuse the type and the size for UUID volumes, therefore,
// buildVolumeMap function will remove them.
func buildVolumeTemplateFromUUID(api *instance.API, volumeUUID string) (*instance.VolumeTemplate, error) {
	res, err := api.GetVolume(&instance.GetVolumeRequest{VolumeID: volumeUUID})
	if err != nil { // FIXME: isNotFoundError
		return nil, fmt.Errorf("Volume %s does not exist.", volumeUUID)
	}

	// Check that volume is not already attached to a server.
	if res.Volume.Server != nil {
		return nil, fmt.Errorf("Volume %s is already attached to %s server.", res.Volume.ID, res.Volume.Server.ID)
	}

	return &instance.VolumeTemplate{
		ID:         res.Volume.ID,
		VolumeType: res.Volume.VolumeType,
		Size:       res.Volume.Size,
	}, nil
}

// validateLocalVolumeSizes validates the total size of local volumes.
func validateLocalVolumeSizes(api *instance.API, volumes []*instance.VolumeTemplate, commercialType string) error {
	// Calculate local volume total size.
	var localVolumeTotalSize scw.Size
	for _, volume := range volumes {
		if volume.VolumeType == instance.VolumeTypeLSSD {
			localVolumeTotalSize += volume.Size
		}
	}

	// Get server types.
	serverTypesRes, err := api.ListServersTypes(&instance.ListServersTypesRequest{})
	if err != nil {
		// Ignore root volume size check.
		logger.Warningf("cannot get server types: %s", err)
		logger.Warningf("skip local volume size validation")
		return nil
	}

	// Validate total size.
	var vConstr *instance.ServerTypeVolumeConstraintSizes
	if st, exists := serverTypesRes.Servers[commercialType]; exists {
		vConstr = st.VolumesConstraint

		if localVolumeTotalSize < vConstr.MinSize || localVolumeTotalSize > vConstr.MaxSize {
			min := humanize.Bytes(uint64(vConstr.MinSize))
			if vConstr.MinSize == vConstr.MaxSize {
				return fmt.Errorf("%s total local volume size must be equal to %s.", commercialType, min)
			}

			max := humanize.Bytes(uint64(vConstr.MaxSize))
			return fmt.Errorf("%s total local volume size must be between %s and %s.", commercialType, min, max)
		}
	} else {
		logger.Warningf("unrecognized server type: %s", commercialType)
		logger.Warningf("skip local volume size validation")
	}

	return nil
}

func validateRootVolume(api *instance.API, image string, rootVolume *instance.VolumeTemplate) error {
	if rootVolume.VolumeType != instance.VolumeTypeLSSD {
		return fmt.Errorf("First volume must be local.")
	}

	if rootVolume.ID != "" {
		// TODO: Improve error message [APIGW-1371]
		return fmt.Errorf("You cannot use an existing volume as a root volume. You must create an image of this volume and use its ID in the 'image' argument.")
	}

	res, err := api.GetImage(&instance.GetImageRequest{ImageID: image})
	if err != nil {
		logger.Warningf("cannot get image %s: %s", image, err)
		logger.Warningf("skip root volume size validation")
	}

	if rootVolume.Size < res.Image.RootVolume.Size {
		return fmt.Errorf("First volume size must be at least %s for this image.", humanize.Bytes(uint64(res.Image.RootVolume.Size)))
	}

	return nil
}

// buildVolumeMap builds the volume map for the server request.
func buildVolumeMap(serverName string, volumes []*instance.VolumeTemplate) map[string]*instance.VolumeTemplate {
	m := make(map[string]*instance.VolumeTemplate)

	for k, v := range volumes {
		idx := strconv.Itoa(k)
		v.Name = serverName + "-" + idx

		// Remove extra data for API validation.
		if v.ID != "" {
			v = &instance.VolumeTemplate{ID: v.ID, Name: v.Name}
		}

		// Remove extra data for API validation.
		if k == 0 {
			v = &instance.VolumeTemplate{Size: v.Size}
		}

		m[idx] = v
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
