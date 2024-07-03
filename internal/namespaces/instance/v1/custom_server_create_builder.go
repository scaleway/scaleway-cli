package instance

import (
	"fmt"
	"net"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

type ServerBuilder struct {
	// createdReq is the request being built
	createReq *instance.CreateServerRequest
	// createIPReq is filled with a request if an IP is needed
	createIPReq *instance.CreateIPRequest

	// All needed APIs
	apiMarketplace *marketplace.API
	apiInstance    *instance.API

	// serverType is filled with the ServerType if CommercialType is found in the API.
	serverType *instance.ServerType
	// serverImage is filled with the Image if one is provided
	serverImage *instance.Image
}

func NewServerBuilder(client *scw.Client, name string, zone scw.Zone, commercialType string) *ServerBuilder {
	sb := &ServerBuilder{
		createReq: &instance.CreateServerRequest{
			Name:           name,
			Zone:           zone,
			CommercialType: commercialType,
		},
		apiMarketplace: marketplace.NewAPI(client),
		apiInstance:    instance.NewAPI(client),
	}

	return sb
}

func (sb *ServerBuilder) AddOrganizationID(orgID *string) *ServerBuilder {
	if orgID != nil {
		sb.createReq.Organization = orgID
	}

	return sb
}

func (sb *ServerBuilder) AddProjectID(projectID *string) *ServerBuilder {
	if projectID != nil {
		sb.createReq.Project = projectID
	}

	return sb
}

func (sb *ServerBuilder) AddEnableIPv6(enableIPv6 *bool) *ServerBuilder {
	sb.createReq.EnableIPv6 = enableIPv6 //nolint: staticcheck

	return sb
}

func (sb *ServerBuilder) AddTags(tags []string) *ServerBuilder {
	if len(tags) > 0 {
		sb.createReq.Tags = tags
	}

	return sb
}

func (sb *ServerBuilder) AddRoutedIPEnabled(routedIPEnabled *bool) *ServerBuilder {
	if routedIPEnabled != nil {
		sb.createReq.RoutedIPEnabled = routedIPEnabled
	}

	return sb
}

func (sb *ServerBuilder) AddAdminPasswordEncryptionSSHKeyID(adminPasswordEncryptionSSHKeyID *string) *ServerBuilder {
	if adminPasswordEncryptionSSHKeyID != nil {
		sb.createReq.AdminPasswordEncryptionSSHKeyID = adminPasswordEncryptionSSHKeyID
	}

	return sb
}

func (sb *ServerBuilder) isWindows() bool {
	return commercialTypeIsWindowsServer(sb.createReq.CommercialType)
}

// defaultIPType returns the default IP type when created by the CLI. Used for ServerBuilder.AddIP
func (sb *ServerBuilder) defaultIPType() instance.IPType {
	if sb.createReq.RoutedIPEnabled != nil {
		if *sb.createReq.RoutedIPEnabled {
			return instance.IPTypeRoutedIPv4
		}
		return instance.IPTypeNat
	}

	return ""
}

// AddImage handle a custom image argument.
// image could be:
//   - A local image UUID.
//   - An image label.
func (sb *ServerBuilder) AddImage(image string) (*ServerBuilder, error) {
	switch {
	case image == "none":
		return sb, nil
	case !validation.IsUUID(image):
		imageLabel := strings.Replace(image, "-", "_", -1)

		localImage, err := sb.apiMarketplace.GetLocalImageByLabel(&marketplace.GetLocalImageByLabelRequest{
			ImageLabel:     imageLabel,
			Zone:           sb.createReq.Zone,
			CommercialType: sb.createReq.CommercialType,
			Type:           marketplace.LocalImageTypeInstanceLocal,
		})
		if err != nil {
			return sb, err
		}

		image = localImage.ID
	}

	sb.createReq.Image = image

	getImageResponse, err := sb.apiInstance.GetImage(&instance.GetImageRequest{
		Zone:    sb.createReq.Zone,
		ImageID: sb.createReq.Image,
	})
	if err != nil {
		logger.Warningf("cannot get image %s: %s", sb.createReq.Image, err)
	} else {
		sb.serverImage = getImageResponse.Image
	}

	sb.serverType = getServerType(sb.apiInstance, sb.createReq.Zone, sb.createReq.CommercialType)

	return sb, nil
}

// AddIP takes an ip argument and change requests accordingly.
// ip could be:
//   - "new"
//   - A flexible IP UUID
//   - A flexible IP address
//   - "dynamic"
//   - "none"
func (sb *ServerBuilder) AddIP(ip string) (*ServerBuilder, error) {
	switch {
	case ip == "" || ip == "new":
		sb.createIPReq = &instance.CreateIPRequest{
			Zone:         sb.createReq.Zone,
			Organization: sb.createReq.Project,
			Project:      sb.createReq.Project,
			Type:         sb.defaultIPType(),
		}
	case validation.IsUUID(ip):
		sb.createReq.PublicIP = scw.StringPtr(ip)
	case net.ParseIP(ip) != nil:
		logger.Debugf("finding public IP UUID from address: %s", ip)
		res, err := sb.apiInstance.GetIP(&instance.GetIPRequest{
			Zone: sb.createReq.Zone,
			IP:   ip,
		})
		if err != nil { // FIXME: isNotFoundError
			return sb, fmt.Errorf("%s does not belong to you", ip)
		}
		sb.createReq.PublicIP = scw.StringPtr(res.IP.ID)

	case ip == "dynamic":
		sb.createReq.DynamicIPRequired = scw.BoolPtr(true)
	case ip == "none":
		sb.createReq.DynamicIPRequired = scw.BoolPtr(false)
	default:
		return sb, fmt.Errorf(`invalid IP "%s", should be either 'new', 'dynamic', 'none', an IP address ID or a reserved flexible IP address`, ip)
	}

	return sb, nil
}

// AddVolumes build volume templates from arguments.
//
// More format details in buildVolumeTemplate function.
//
// Also add default volumes to server, ex: scratch storage for GPU servers
func (sb *ServerBuilder) AddVolumes(rootVolume string, additionalVolumes []string) (*ServerBuilder, error) {
	if len(additionalVolumes) > 0 || rootVolume != "" {
		// Create initial volume template map.
		volumes, err := buildVolumes(sb.apiInstance, sb.createReq.Zone, sb.createReq.Name, rootVolume, additionalVolumes)
		if err != nil {
			return sb, err
		}

		// Validate root volume type and size.
		if sb.serverImage != nil {
			if err := validateRootVolume(sb.serverImage.RootVolume.Size, volumes["0"]); err != nil {
				return sb, err
			}
		} else {
			logger.Warningf("skipping root volume validation")
		}

		// Validate total local volume sizes.
		if sb.serverType != nil && sb.serverImage != nil {
			if err := validateLocalVolumeSizes(volumes, sb.serverType, sb.createReq.CommercialType, sb.serverImage.RootVolume.Size); err != nil {
				return sb, err
			}
		} else {
			logger.Warningf("skip local volume size validation")
		}

		// Sanitize the volume map to respect API schemas
		sb.createReq.Volumes = sanitizeVolumeMap(sb.createReq.Name, volumes)
	}

	if sb.serverType != nil {
		sb.createReq.Volumes = addDefaultVolumes(sb.serverType, sb.createReq.Volumes)
	}

	return sb, nil
}

func (sb *ServerBuilder) AddBootType(bootType string) *ServerBuilder {
	instanceBootType := instance.BootType(bootType)
	sb.createReq.BootType = &instanceBootType

	return sb
}

func (sb *ServerBuilder) AddSecurityGroup(securityGroupID string) *ServerBuilder {
	if securityGroupID != "" {
		sb.createReq.SecurityGroup = &securityGroupID
	}

	return sb
}

func (sb *ServerBuilder) AddPlacementGroup(placementGroupID string) *ServerBuilder {
	if placementGroupID != "" {
		sb.createReq.PlacementGroup = &placementGroupID
	}

	return sb
}

func (sb *ServerBuilder) Validate() error {
	if sb.isWindows() && sb.createReq.AdminPasswordEncryptionSSHKeyID == nil {
		return &core.CliError{
			Err:     core.MissingRequiredArgumentError("admin-password-encryption-ssh-key-id").Err,
			Details: "Expected a SSH Key ID to encrypt Admin RDP password. If not provided, no password will be generated. Key must be RSA Public Key.",
			Hint:    "Use completion or get your ssh key id using 'scw iam ssh-key list',",
			Code:    1,
			Empty:   false,
		}
	}

	if sb.serverType != nil && sb.serverImage != nil {
		if err := validateImageServerTypeCompatibility(sb.serverImage, sb.serverType, sb.createReq.CommercialType); err != nil {
			return err
		}
	} else {
		logger.Warningf("skipping image server-type compatibility validation")
	}

	return nil
}

func (sb *ServerBuilder) Build() (*instance.CreateServerRequest, *instance.CreateIPRequest) {
	return sb.createReq, sb.createIPReq
}
