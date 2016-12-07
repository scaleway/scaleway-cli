// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Interact with Scaleway API

// Package api contains client and functions to interact with Scaleway API
package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"golang.org/x/sync/errgroup"
)

// Default values
var (
	AccountAPI     = "https://account.scaleway.com/"
	MetadataAPI    = "http://169.254.42.42/"
	MarketplaceAPI = "https://api-marketplace.scaleway.com"
	ComputeAPIPar1 = "https://cp-par1.scaleway.com/"
	ComputeAPIAms1 = "https://cp-ams1.scaleway.com"

	URLPublicDNS  = ".pub.cloud.scaleway.com"
	URLPrivateDNS = ".priv.cloud.scaleway.com"
)

func init() {
	if url := os.Getenv("SCW_ACCOUNT_API"); url != "" {
		AccountAPI = url
	}
	if url := os.Getenv("SCW_METADATA_API"); url != "" {
		MetadataAPI = url
	}
	if url := os.Getenv("SCW_MARKETPLACE_API"); url != "" {
		MarketplaceAPI = url
	}
}

const (
	perPage = 50
)

// ScalewayAPI is the interface used to communicate with the Scaleway API
type ScalewayAPI struct {
	// Organization is the identifier of the Scaleway organization
	Organization string

	// Token is the authentication token for the Scaleway organization
	Token string

	// Password is the authentication password
	password string

	userAgent string

	// Cache is used to quickly resolve identifiers from names
	Cache *ScalewayCache

	client     *http.Client
	verbose    bool
	computeAPI string

	Region string
	//
	Logger
}

// ScalewayAPIError represents a Scaleway API Error
type ScalewayAPIError struct {
	// Message is a human-friendly error message
	APIMessage string `json:"message,omitempty"`

	// Type is a string code that defines the kind of error
	Type string `json:"type,omitempty"`

	// Fields contains detail about validation error
	Fields map[string][]string `json:"fields,omitempty"`

	// StatusCode is the HTTP status code received
	StatusCode int `json:"-"`

	// Message
	Message string `json:"-"`
}

// Error returns a string representing the error
func (e ScalewayAPIError) Error() string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "StatusCode: %v, ", e.StatusCode)
	fmt.Fprintf(&b, "Type: %v, ", e.Type)
	fmt.Fprintf(&b, "APIMessage: \x1b[31m%v\x1b[0m", e.APIMessage)
	if len(e.Fields) > 0 {
		fmt.Fprintf(&b, ", Details: %v", e.Fields)
	}
	return b.String()
}

// HideAPICredentials removes API credentials from a string
func (s *ScalewayAPI) HideAPICredentials(input string) string {
	output := input
	if s.Token != "" {
		output = strings.Replace(output, s.Token, "00000000-0000-4000-8000-000000000000", -1)
	}
	if s.Organization != "" {
		output = strings.Replace(output, s.Organization, "00000000-0000-5000-9000-000000000000", -1)
	}
	if s.password != "" {
		output = strings.Replace(output, s.password, "XX-XX-XX-XX", -1)
	}
	return output
}

// ScalewayIPAddress represents a Scaleway IP address
type ScalewayIPAddress struct {
	// Identifier is a unique identifier for the IP address
	Identifier string `json:"id,omitempty"`

	// IP is an IPv4 address
	IP string `json:"address,omitempty"`

	// Dynamic is a flag that defines an IP that change on each reboot
	Dynamic *bool `json:"dynamic,omitempty"`
}

// ScalewayVolume represents a Scaleway Volume
type ScalewayVolume struct {
	// Identifier is a unique identifier for the volume
	Identifier string `json:"id,omitempty"`

	// Size is the allocated size of the volume
	Size uint64 `json:"size,omitempty"`

	// CreationDate is the creation date of the volume
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the volume
	ModificationDate string `json:"modification_date,omitempty"`

	// Organization is the organization owning the volume
	Organization string `json:"organization,omitempty"`

	// Name is the name of the volume
	Name string `json:"name,omitempty"`

	// Server is the server using this image
	Server *struct {
		Identifier string `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"server,omitempty"`

	// VolumeType is a Scaleway identifier for the kind of volume (default: l_ssd)
	VolumeType string `json:"volume_type,omitempty"`

	// ExportURI represents the url used by initrd/scripts to attach the volume
	ExportURI string `json:"export_uri,omitempty"`
}

// ScalewayOneVolume represents the response of a GET /volumes/UUID API call
type ScalewayOneVolume struct {
	Volume ScalewayVolume `json:"volume,omitempty"`
}

// ScalewayVolumes represents a group of Scaleway volumes
type ScalewayVolumes struct {
	// Volumes holds scaleway volumes of the response
	Volumes []ScalewayVolume `json:"volumes,omitempty"`
}

// ScalewayVolumeDefinition represents a Scaleway volume definition
type ScalewayVolumeDefinition struct {
	// Name is the user-defined name of the volume
	Name string `json:"name"`

	// Image is the image used by the volume
	Size uint64 `json:"size"`

	// Bootscript is the bootscript used by the volume
	Type string `json:"volume_type"`

	// Organization is the owner of the volume
	Organization string `json:"organization"`
}

// ScalewayVolumePutDefinition represents a Scaleway volume with nullable fields (for PUT)
type ScalewayVolumePutDefinition struct {
	Identifier       *string `json:"id,omitempty"`
	Size             *uint64 `json:"size,omitempty"`
	CreationDate     *string `json:"creation_date,omitempty"`
	ModificationDate *string `json:"modification_date,omitempty"`
	Organization     *string `json:"organization,omitempty"`
	Name             *string `json:"name,omitempty"`
	Server           struct {
		Identifier *string `json:"id,omitempty"`
		Name       *string `json:"name,omitempty"`
	} `json:"server,omitempty"`
	VolumeType *string `json:"volume_type,omitempty"`
	ExportURI  *string `json:"export_uri,omitempty"`
}

// ScalewayImage represents a Scaleway Image
type ScalewayImage struct {
	// Identifier is a unique identifier for the image
	Identifier string `json:"id,omitempty"`

	// Name is a user-defined name for the image
	Name string `json:"name,omitempty"`

	// CreationDate is the creation date of the image
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the image
	ModificationDate string `json:"modification_date,omitempty"`

	// RootVolume is the root volume bound to the image
	RootVolume ScalewayVolume `json:"root_volume,omitempty"`

	// Public is true for public images and false for user images
	Public bool `json:"public,omitempty"`

	// Bootscript is the bootscript bound to the image
	DefaultBootscript *ScalewayBootscript `json:"default_bootscript,omitempty"`

	// Organization is the owner of the image
	Organization string `json:"organization,omitempty"`

	// Arch is the architecture target of the image
	Arch string `json:"arch,omitempty"`

	// FIXME: extra_volumes
}

// ScalewayImageIdentifier represents a Scaleway Image Identifier
type ScalewayImageIdentifier struct {
	Identifier string
	Arch       string
	Region     string
	Owner      string
}

// ScalewayOneImage represents the response of a GET /images/UUID API call
type ScalewayOneImage struct {
	Image ScalewayImage `json:"image,omitempty"`
}

// ScalewayImages represents a group of Scaleway images
type ScalewayImages struct {
	// Images holds scaleway images of the response
	Images []ScalewayImage `json:"images,omitempty"`
}

// ScalewaySnapshot represents a Scaleway Snapshot
type ScalewaySnapshot struct {
	// Identifier is a unique identifier for the snapshot
	Identifier string `json:"id,omitempty"`

	// Name is a user-defined name for the snapshot
	Name string `json:"name,omitempty"`

	// CreationDate is the creation date of the snapshot
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the snapshot
	ModificationDate string `json:"modification_date,omitempty"`

	// Size is the allocated size of the volume
	Size uint64 `json:"size,omitempty"`

	// Organization is the owner of the snapshot
	Organization string `json:"organization"`

	// State is the current state of the snapshot
	State string `json:"state"`

	// VolumeType is the kind of volume behind the snapshot
	VolumeType string `json:"volume_type"`

	// BaseVolume is the volume from which the snapshot inherits
	BaseVolume ScalewayVolume `json:"base_volume,omitempty"`
}

// ScalewayOneSnapshot represents the response of a GET /snapshots/UUID API call
type ScalewayOneSnapshot struct {
	Snapshot ScalewaySnapshot `json:"snapshot,omitempty"`
}

// ScalewaySnapshots represents a group of Scaleway snapshots
type ScalewaySnapshots struct {
	// Snapshots holds scaleway snapshots of the response
	Snapshots []ScalewaySnapshot `json:"snapshots,omitempty"`
}

// ScalewayBootscript represents a Scaleway Bootscript
type ScalewayBootscript struct {
	Bootcmdargs string `json:"bootcmdargs,omitempty"`
	Dtb         string `json:"dtb,omitempty"`
	Initrd      string `json:"initrd,omitempty"`
	Kernel      string `json:"kernel,omitempty"`

	// Arch is the architecture target of the bootscript
	Arch string `json:"architecture,omitempty"`

	// Identifier is a unique identifier for the bootscript
	Identifier string `json:"id,omitempty"`

	// Organization is the owner of the bootscript
	Organization string `json:"organization,omitempty"`

	// Name is a user-defined name for the bootscript
	Title string `json:"title,omitempty"`

	// Public is true for public bootscripts and false for user bootscripts
	Public bool `json:"public,omitempty"`

	Default bool `json:"default,omitempty"`
}

// ScalewayOneBootscript represents the response of a GET /bootscripts/UUID API call
type ScalewayOneBootscript struct {
	Bootscript ScalewayBootscript `json:"bootscript,omitempty"`
}

// ScalewayBootscripts represents a group of Scaleway bootscripts
type ScalewayBootscripts struct {
	// Bootscripts holds Scaleway bootscripts of the response
	Bootscripts []ScalewayBootscript `json:"bootscripts,omitempty"`
}

// ScalewayTask represents a Scaleway Task
type ScalewayTask struct {
	// Identifier is a unique identifier for the task
	Identifier string `json:"id,omitempty"`

	// StartDate is the start date of the task
	StartDate string `json:"started_at,omitempty"`

	// TerminationDate is the termination date of the task
	TerminationDate string `json:"terminated_at,omitempty"`

	HrefFrom string `json:"href_from,omitempty"`

	Description string `json:"description,omitempty"`

	Status string `json:"status,omitempty"`

	Progress int `json:"progress,omitempty"`
}

// ScalewayOneTask represents the response of a GET /tasks/UUID API call
type ScalewayOneTask struct {
	Task ScalewayTask `json:"task,omitempty"`
}

// ScalewayTasks represents a group of Scaleway tasks
type ScalewayTasks struct {
	// Tasks holds scaleway tasks of the response
	Tasks []ScalewayTask `json:"tasks,omitempty"`
}

// ScalewaySecurityGroupRule definition
type ScalewaySecurityGroupRule struct {
	Direction    string `json:"direction"`
	Protocol     string `json:"protocol"`
	IPRange      string `json:"ip_range"`
	DestPortFrom int    `json:"dest_port_from,omitempty"`
	Action       string `json:"action"`
	Position     int    `json:"position"`
	DestPortTo   string `json:"dest_port_to"`
	Editable     bool   `json:"editable"`
	ID           string `json:"id"`
}

// ScalewayGetSecurityGroupRules represents the response of a GET /security_group/{groupID}/rules
type ScalewayGetSecurityGroupRules struct {
	Rules []ScalewaySecurityGroupRule `json:"rules"`
}

// ScalewayGetSecurityGroupRule represents the response of a GET /security_group/{groupID}/rules/{ruleID}
type ScalewayGetSecurityGroupRule struct {
	Rules ScalewaySecurityGroupRule `json:"rule"`
}

// ScalewayNewSecurityGroupRule definition POST/PUT request /security_group/{groupID}
type ScalewayNewSecurityGroupRule struct {
	Action       string `json:"action"`
	Direction    string `json:"direction"`
	IPRange      string `json:"ip_range"`
	Protocol     string `json:"protocol"`
	DestPortFrom int    `json:"dest_port_from,omitempty"`
}

// ScalewaySecurityGroups definition
type ScalewaySecurityGroups struct {
	Description           string                  `json:"description"`
	ID                    string                  `json:"id"`
	Organization          string                  `json:"organization"`
	Name                  string                  `json:"name"`
	Servers               []ScalewaySecurityGroup `json:"servers"`
	EnableDefaultSecurity bool                    `json:"enable_default_security"`
	OrganizationDefault   bool                    `json:"organization_default"`
}

// ScalewayGetSecurityGroups represents the response of a GET /security_groups/
type ScalewayGetSecurityGroups struct {
	SecurityGroups []ScalewaySecurityGroups `json:"security_groups"`
}

// ScalewayGetSecurityGroup represents the response of a GET /security_groups/{groupID}
type ScalewayGetSecurityGroup struct {
	SecurityGroups ScalewaySecurityGroups `json:"security_group"`
}

// ScalewayIPDefinition represents the IP's fields
type ScalewayIPDefinition struct {
	Organization string  `json:"organization"`
	Reverse      *string `json:"reverse"`
	ID           string  `json:"id"`
	Server       *struct {
		Identifier string `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"server"`
	Address string `json:"address"`
}

// ScalewayGetIPS represents the response of a GET /ips/
type ScalewayGetIPS struct {
	IPS []ScalewayIPDefinition `json:"ips"`
}

// ScalewayGetIP represents the response of a GET /ips/{id_ip}
type ScalewayGetIP struct {
	IP ScalewayIPDefinition `json:"ip"`
}

// ScalewaySecurityGroup represents a Scaleway security group
type ScalewaySecurityGroup struct {
	// Identifier is a unique identifier for the security group
	Identifier string `json:"id,omitempty"`

	// Name is the user-defined name of the security group
	Name string `json:"name,omitempty"`
}

// ScalewayNewSecurityGroup definition POST request /security_groups
type ScalewayNewSecurityGroup struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

// ScalewayUpdateSecurityGroup definition PUT request /security_groups
type ScalewayUpdateSecurityGroup struct {
	Organization        string `json:"organization"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	OrganizationDefault bool   `json:"organization_default"`
}

// ScalewayServer represents a Scaleway server
type ScalewayServer struct {
	// Arch is the architecture target of the server
	Arch string `json:"arch,omitempty"`

	// Identifier is a unique identifier for the server
	Identifier string `json:"id,omitempty"`

	// Name is the user-defined name of the server
	Name string `json:"name,omitempty"`

	// CreationDate is the creation date of the server
	CreationDate string `json:"creation_date,omitempty"`

	// ModificationDate is the date of the last modification of the server
	ModificationDate string `json:"modification_date,omitempty"`

	// Image is the image used by the server
	Image ScalewayImage `json:"image,omitempty"`

	// DynamicIPRequired is a flag that defines a server with a dynamic ip address attached
	DynamicIPRequired *bool `json:"dynamic_ip_required,omitempty"`

	// PublicIP is the public IP address bound to the server
	PublicAddress ScalewayIPAddress `json:"public_ip,omitempty"`

	// State is the current status of the server
	State string `json:"state,omitempty"`

	// StateDetail is the detailed status of the server
	StateDetail string `json:"state_detail,omitempty"`

	// PrivateIP represents the private IPV4 attached to the server (changes on each boot)
	PrivateIP string `json:"private_ip,omitempty"`

	// Bootscript is the unique identifier of the selected bootscript
	Bootscript *ScalewayBootscript `json:"bootscript,omitempty"`

	// Hostname represents the ServerName in a format compatible with unix's hostname
	Hostname string `json:"hostname,omitempty"`

	// Tags represents user-defined tags
	Tags []string `json:"tags,omitempty"`

	// Volumes are the attached volumes
	Volumes map[string]ScalewayVolume `json:"volumes,omitempty"`

	// SecurityGroup is the selected security group object
	SecurityGroup ScalewaySecurityGroup `json:"security_group,omitempty"`

	// Organization is the owner of the server
	Organization string `json:"organization,omitempty"`

	// CommercialType is the commercial type of the server (i.e: C1, C2[SML], VC1S)
	CommercialType string `json:"commercial_type,omitempty"`

	// Location of the server
	Location struct {
		Platform   string `json:"platform_id,omitempty"`
		Chassis    string `json:"chassis_id,omitempty"`
		Cluster    string `json:"cluster_id,omitempty"`
		Hypervisor string `json:"hypervisor_id,omitempty"`
		Blade      string `json:"blade_id,omitempty"`
		Node       string `json:"node_id,omitempty"`
		ZoneID     string `json:"zone_id,omitempty"`
	} `json:"location,omitempty"`

	IPV6 *ScalewayIPV6Definition `json:"ipv6,omitempty"`

	EnableIPV6 bool `json:"enable_ipv6,omitempty"`

	// This fields are not returned by the API, we generate it
	DNSPublic  string `json:"dns_public,omitempty"`
	DNSPrivate string `json:"dns_private,omitempty"`
}

// ScalewayIPV6Definition represents a Scaleway ipv6
type ScalewayIPV6Definition struct {
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Address string `json:"address"`
}

// ScalewayServerPatchDefinition represents a Scaleway server with nullable fields (for PATCH)
type ScalewayServerPatchDefinition struct {
	Arch              *string                    `json:"arch,omitempty"`
	Name              *string                    `json:"name,omitempty"`
	CreationDate      *string                    `json:"creation_date,omitempty"`
	ModificationDate  *string                    `json:"modification_date,omitempty"`
	Image             *ScalewayImage             `json:"image,omitempty"`
	DynamicIPRequired *bool                      `json:"dynamic_ip_required,omitempty"`
	PublicAddress     *ScalewayIPAddress         `json:"public_ip,omitempty"`
	State             *string                    `json:"state,omitempty"`
	StateDetail       *string                    `json:"state_detail,omitempty"`
	PrivateIP         *string                    `json:"private_ip,omitempty"`
	Bootscript        *string                    `json:"bootscript,omitempty"`
	Hostname          *string                    `json:"hostname,omitempty"`
	Volumes           *map[string]ScalewayVolume `json:"volumes,omitempty"`
	SecurityGroup     *ScalewaySecurityGroup     `json:"security_group,omitempty"`
	Organization      *string                    `json:"organization,omitempty"`
	Tags              *[]string                  `json:"tags,omitempty"`
	IPV6              *ScalewayIPV6Definition    `json:"ipv6,omitempty"`
	EnableIPV6        *bool                      `json:"enable_ipv6,omitempty"`
}

// ScalewayServerDefinition represents a Scaleway server with image definition
type ScalewayServerDefinition struct {
	// Name is the user-defined name of the server
	Name string `json:"name"`

	// Image is the image used by the server
	Image *string `json:"image,omitempty"`

	// Volumes are the attached volumes
	Volumes map[string]string `json:"volumes,omitempty"`

	// DynamicIPRequired is a flag that defines a server with a dynamic ip address attached
	DynamicIPRequired *bool `json:"dynamic_ip_required,omitempty"`

	// Bootscript is the bootscript used by the server
	Bootscript *string `json:"bootscript"`

	// Tags are the metadata tags attached to the server
	Tags []string `json:"tags,omitempty"`

	// Organization is the owner of the server
	Organization string `json:"organization"`

	// CommercialType is the commercial type of the server (i.e: C1, C2[SML], VC1S)
	CommercialType string `json:"commercial_type"`

	PublicIP string `json:"public_ip,omitempty"`

	EnableIPV6 bool `json:"enable_ipv6,omitempty"`

	SecurityGroup string `json:"security_group,omitempty"`
}

// ScalewayOneServer represents the response of a GET /servers/UUID API call
type ScalewayOneServer struct {
	Server ScalewayServer `json:"server,omitempty"`
}

// ScalewayServers represents a group of Scaleway servers
type ScalewayServers struct {
	// Servers holds scaleway servers of the response
	Servers []ScalewayServer `json:"servers,omitempty"`
}

// ScalewayServerAction represents an action to perform on a Scaleway server
type ScalewayServerAction struct {
	// Action is the name of the action to trigger
	Action string `json:"action,omitempty"`
}

// ScalewaySnapshotDefinition represents a Scaleway snapshot definition
type ScalewaySnapshotDefinition struct {
	VolumeIDentifier string `json:"volume_id"`
	Name             string `json:"name,omitempty"`
	Organization     string `json:"organization"`
}

// ScalewayImageDefinition represents a Scaleway image definition
type ScalewayImageDefinition struct {
	SnapshotIDentifier string  `json:"root_volume"`
	Name               string  `json:"name,omitempty"`
	Organization       string  `json:"organization"`
	Arch               string  `json:"arch"`
	DefaultBootscript  *string `json:"default_bootscript,omitempty"`
}

// ScalewayRoleDefinition represents a Scaleway Token UserId Role
type ScalewayRoleDefinition struct {
	Organization ScalewayOrganizationDefinition `json:"organization,omitempty"`
	Role         string                         `json:"role,omitempty"`
}

// ScalewayTokenDefinition represents a Scaleway Token
type ScalewayTokenDefinition struct {
	UserID             string                 `json:"user_id"`
	Description        string                 `json:"description,omitempty"`
	Roles              ScalewayRoleDefinition `json:"roles"`
	Expires            string                 `json:"expires"`
	InheritsUsersPerms bool                   `json:"inherits_user_perms"`
	ID                 string                 `json:"id"`
}

// ScalewayTokensDefinition represents a Scaleway Tokens
type ScalewayTokensDefinition struct {
	Token ScalewayTokenDefinition `json:"token"`
}

// ScalewayGetTokens represents a list of Scaleway Tokens
type ScalewayGetTokens struct {
	Tokens []ScalewayTokenDefinition `json:"tokens"`
}

// ScalewayContainerData represents a Scaleway container data (S3)
type ScalewayContainerData struct {
	LastModified string `json:"last_modified"`
	Name         string `json:"name"`
	Size         string `json:"size"`
}

// ScalewayGetContainerDatas represents a list of Scaleway containers data (S3)
type ScalewayGetContainerDatas struct {
	Container []ScalewayContainerData `json:"container"`
}

// ScalewayContainer represents a Scaleway container (S3)
type ScalewayContainer struct {
	ScalewayOrganizationDefinition `json:"organization"`
	Name                           string `json:"name"`
	Size                           string `json:"size"`
}

// ScalewayGetContainers represents a list of Scaleway containers (S3)
type ScalewayGetContainers struct {
	Containers []ScalewayContainer `json:"containers"`
}

// ScalewayConnectResponse represents the answer from POST /tokens
type ScalewayConnectResponse struct {
	Token ScalewayTokenDefinition `json:"token"`
}

// ScalewayConnect represents the data to connect
type ScalewayConnect struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Description string `json:"description"`
	Expires     bool   `json:"expires"`
}

// ScalewayOrganizationDefinition represents a Scaleway Organization
type ScalewayOrganizationDefinition struct {
	ID    string                   `json:"id"`
	Name  string                   `json:"name"`
	Users []ScalewayUserDefinition `json:"users"`
}

// ScalewayOrganizationsDefinition represents a Scaleway Organizations
type ScalewayOrganizationsDefinition struct {
	Organizations []ScalewayOrganizationDefinition `json:"organizations"`
}

// ScalewayUserDefinition represents a Scaleway User
type ScalewayUserDefinition struct {
	Email         string                           `json:"email"`
	Firstname     string                           `json:"firstname"`
	Fullname      string                           `json:"fullname"`
	ID            string                           `json:"id"`
	Lastname      string                           `json:"lastname"`
	Organizations []ScalewayOrganizationDefinition `json:"organizations"`
	Roles         []ScalewayRoleDefinition         `json:"roles"`
	SSHPublicKeys []ScalewayKeyDefinition          `json:"ssh_public_keys"`
}

// ScalewayUsersDefinition represents the response of a GET /user
type ScalewayUsersDefinition struct {
	User ScalewayUserDefinition `json:"user"`
}

// ScalewayKeyDefinition represents a key
type ScalewayKeyDefinition struct {
	Key         string `json:"key"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// ScalewayUserPatchSSHKeyDefinition represents a User Patch
type ScalewayUserPatchSSHKeyDefinition struct {
	SSHPublicKeys []ScalewayKeyDefinition `json:"ssh_public_keys"`
}

// ScalewayDashboardResp represents a dashboard received from the API
type ScalewayDashboardResp struct {
	Dashboard ScalewayDashboard
}

// ScalewayDashboard represents a dashboard
type ScalewayDashboard struct {
	VolumesCount        int `json:"volumes_count"`
	RunningServersCount int `json:"running_servers_count"`
	ImagesCount         int `json:"images_count"`
	SnapshotsCount      int `json:"snapshots_count"`
	ServersCount        int `json:"servers_count"`
	IPsCount            int `json:"ips_count"`
}

// ScalewayPermissions represents the response of GET /permissions
type ScalewayPermissions map[string]ScalewayPermCategory

// ScalewayPermCategory represents ScalewayPermissions's fields
type ScalewayPermCategory map[string][]string

// ScalewayPermissionDefinition represents the permissions
type ScalewayPermissionDefinition struct {
	Permissions ScalewayPermissions `json:"permissions"`
}

// ScalewayUserdatas represents the response of a GET /user_data
type ScalewayUserdatas struct {
	UserData []string `json:"user_data"`
}

// ScalewayQuota represents a map of quota (name, value)
type ScalewayQuota map[string]int

// ScalewayGetQuotas represents the response of GET /organizations/{orga_id}/quotas
type ScalewayGetQuotas struct {
	Quotas ScalewayQuota `json:"quotas"`
}

// ScalewayUserdata represents []byte
type ScalewayUserdata []byte

// FuncMap used for json inspection
var FuncMap = template.FuncMap{
	"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
}

// MarketLocalImageDefinition represents localImage of marketplace version
type MarketLocalImageDefinition struct {
	Arch string `json:"arch"`
	ID   string `json:"id"`
	Zone string `json:"zone"`
}

// MarketLocalImages represents an array of local images
type MarketLocalImages struct {
	LocalImages []MarketLocalImageDefinition `json:"local_images"`
}

// MarketLocalImage represents local image
type MarketLocalImage struct {
	LocalImages MarketLocalImageDefinition `json:"local_image"`
}

// MarketVersionDefinition represents version of marketplace image
type MarketVersionDefinition struct {
	CreationDate string `json:"creation_date"`
	ID           string `json:"id"`
	Image        struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"image"`
	ModificationDate string `json:"modification_date"`
	Name             string `json:"name"`
	MarketLocalImages
}

// MarketVersions represents an array of marketplace image versions
type MarketVersions struct {
	Versions []MarketVersionDefinition `json:"versions"`
}

// MarketVersion represents version of marketplace image
type MarketVersion struct {
	Version MarketVersionDefinition `json:"version"`
}

// MarketImage represents MarketPlace image
type MarketImage struct {
	Categories           []string `json:"categories"`
	CreationDate         string   `json:"creation_date"`
	CurrentPublicVersion string   `json:"current_public_version"`
	Description          string   `json:"description"`
	ID                   string   `json:"id"`
	Logo                 string   `json:"logo"`
	ModificationDate     string   `json:"modification_date"`
	Name                 string   `json:"name"`
	Organization         struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"organization"`
	Public bool `json:"-"`
	MarketVersions
}

// MarketImages represents MarketPlace images
type MarketImages struct {
	Images []MarketImage `json:"images"`
}

// NewScalewayAPI creates a ready-to-use ScalewayAPI client
func NewScalewayAPI(organization, token, userAgent, region string, options ...func(*ScalewayAPI)) (*ScalewayAPI, error) {
	s := &ScalewayAPI{
		// exposed
		Organization: organization,
		Token:        token,
		Logger:       NewDefaultLogger(),

		// internal
		client:    &http.Client{},
		verbose:   os.Getenv("SCW_VERBOSE_API") != "",
		password:  "",
		userAgent: userAgent,
	}
	for _, option := range options {
		option(s)
	}
	cache, err := NewScalewayCache(func() { s.Logger.Debugf("Writing cache file to disk") })
	if err != nil {
		return nil, err
	}
	s.Cache = cache
	if os.Getenv("SCW_TLSVERIFY") == "0" {
		s.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	switch region {
	case "par1", "":
		s.computeAPI = ComputeAPIPar1
	case "ams1":
		s.computeAPI = ComputeAPIAms1
	default:
		return nil, fmt.Errorf("%s isn't a valid region", region)
	}
	s.Region = region
	if url := os.Getenv("SCW_COMPUTE_API"); url != "" {
		s.computeAPI = url
	}
	return s, nil
}

// ClearCache clears the cache
func (s *ScalewayAPI) ClearCache() {
	s.Cache.Clear()
}

// Sync flushes out the cache to the disk
func (s *ScalewayAPI) Sync() {
	s.Cache.Save()
}

func (s *ScalewayAPI) response(method, uri string, content io.Reader) (resp *http.Response, err error) {
	var (
		req *http.Request
	)

	req, err = http.NewRequest(method, uri, content)
	if err != nil {
		err = fmt.Errorf("response %s %s", method, uri)
		return
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", s.userAgent)
	s.LogHTTP(req)
	if s.verbose {
		dump, _ := httputil.DumpRequest(req, true)
		s.Debugf("%v", string(dump))
	} else {
		s.Debugf("[%s]: %v", method, uri)
	}
	resp, err = s.client.Do(req)
	return
}

// GetResponsePaginate fetchs all resources and returns an http.Response object for the requested resource
func (s *ScalewayAPI) GetResponsePaginate(apiURL, resource string, values url.Values) (*http.Response, error) {
	resp, err := s.response("HEAD", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, values.Encode()), nil)
	if err != nil {
		return nil, err
	}

	count := resp.Header.Get("X-Total-Count")
	var maxElem int
	if count == "" {
		maxElem = 0
	} else {
		maxElem, err = strconv.Atoi(count)
		if err != nil {
			return nil, err
		}
	}

	get := maxElem / perPage
	if (float32(maxElem) / perPage) > float32(get) {
		get++
	}

	if get <= 1 { // If there is 0 or 1 page of result, the response is not paginated
		if len(values) == 0 {
			return s.response("GET", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), nil)
		}
		return s.response("GET", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, values.Encode()), nil)
	}

	fetchAll := !(values.Get("per_page") != "" || values.Get("page") != "")
	if fetchAll {
		var g errgroup.Group

		ch := make(chan *http.Response, get)
		for i := 1; i <= get; i++ {
			i := i // closure tricks
			g.Go(func() (err error) {
				var resp *http.Response

				val := url.Values{}
				val.Set("per_page", fmt.Sprintf("%v", perPage))
				val.Set("page", fmt.Sprintf("%v", i))
				resp, err = s.response("GET", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, val.Encode()), nil)
				ch <- resp
				return
			})
		}
		if err = g.Wait(); err != nil {
			return nil, err
		}
		newBody := make(map[string][]json.RawMessage)
		body := make(map[string][]json.RawMessage)
		key := ""
		for i := 0; i < get; i++ {
			res := <-ch
			if res.StatusCode != http.StatusOK {
				return res, nil
			}
			content, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(content, &body); err != nil {
				return nil, err
			}

			if i == 0 {
				resp = res
				for k := range body {
					key = k
					break
				}
			}
			newBody[key] = append(newBody[key], body[key]...)
		}
		payload := new(bytes.Buffer)
		if err := json.NewEncoder(payload).Encode(newBody); err != nil {
			return nil, err
		}
		resp.Body = ioutil.NopCloser(payload)
	} else {
		resp, err = s.response("GET", fmt.Sprintf("%s/%s?%s", strings.TrimRight(apiURL, "/"), resource, values.Encode()), nil)
	}
	return resp, err
}

// PostResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PostResponse(apiURL, resource string, data interface{}) (*http.Response, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}
	return s.response("POST", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), payload)
}

// PatchResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PatchResponse(apiURL, resource string, data interface{}) (*http.Response, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}
	return s.response("PATCH", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), payload)
}

// PutResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PutResponse(apiURL, resource string, data interface{}) (*http.Response, error) {
	payload := new(bytes.Buffer)
	if err := json.NewEncoder(payload).Encode(data); err != nil {
		return nil, err
	}
	return s.response("PUT", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), payload)
}

// DeleteResponse returns an http.Response object for the deleted resource
func (s *ScalewayAPI) DeleteResponse(apiURL, resource string) (*http.Response, error) {
	return s.response("DELETE", fmt.Sprintf("%s/%s", strings.TrimRight(apiURL, "/"), resource), nil)
}

// handleHTTPError checks the statusCode and displays the error
func (s *ScalewayAPI) handleHTTPError(goodStatusCode []int, resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if s.verbose {
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			var js bytes.Buffer

			err = json.Indent(&js, body, "", "  ")
			if err != nil {
				s.Debugf("[Response]: [%v]\n%v", resp.StatusCode, string(dump))
			} else {
				s.Debugf("[Response]: [%v]\n%v", resp.StatusCode, js.String())
			}
		}
	} else {
		s.Debugf("[Response]: [%v]\n%v", resp.StatusCode, string(body))
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, errors.New(string(body))
	}
	good := false
	for _, code := range goodStatusCode {
		if code == resp.StatusCode {
			good = true
		}
	}
	if !good {
		var scwError ScalewayAPIError

		if err := json.Unmarshal(body, &scwError); err != nil {
			return nil, err
		}
		scwError.StatusCode = resp.StatusCode
		s.Debugf("%s", scwError.Error())
		return nil, scwError
	}
	return body, nil
}

// ScalewaySortServers represents a wrapper to sort by CreationDate the servers
type ScalewaySortServers []ScalewayServer

func (s ScalewaySortServers) Len() int {
	return len(s)
}

func (s ScalewaySortServers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ScalewaySortServers) Less(i, j int) bool {
	date1, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", s[i].CreationDate)
	date2, _ := time.Parse("2006-01-02T15:04:05.000000+00:00", s[j].CreationDate)
	return date2.Before(date1)
}

// GetTasks get the list of tasks from the ScalewayAPI
func (s *ScalewayAPI) GetTasks() (*[]ScalewayTask, error) {
	query := url.Values{}
	resp, err := s.GetResponsePaginate(s.computeAPI, "tasks", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var tasks ScalewayTasks

	if err = json.Unmarshal(body, &tasks); err != nil {
		return nil, err
	}
	return &tasks.Tasks, nil
}

// CheckCredentials performs a dummy check to ensure we can contact the API
func (s *ScalewayAPI) CheckCredentials() error {
	query := url.Values{}

	resp, err := s.GetResponsePaginate(AccountAPI, "tokens", query)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return err
	}
	found := false
	var tokens ScalewayGetTokens

	if err = json.Unmarshal(body, &tokens); err != nil {
		return err
	}
	for _, token := range tokens.Tokens {
		if token.ID == s.Token {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Invalid token %v", s.Token)
	}
	return nil
}

// GetUserID returns the userID
func (s *ScalewayAPI) GetUserID() (string, error) {
	resp, err := s.GetResponsePaginate(AccountAPI, fmt.Sprintf("tokens/%s", s.Token), url.Values{})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return "", err
	}
	var token ScalewayTokensDefinition

	if err = json.Unmarshal(body, &token); err != nil {
		return "", err
	}
	return token.Token.UserID, nil
}

// GetOrganization returns Organization
func (s *ScalewayAPI) GetOrganization() (*ScalewayOrganizationsDefinition, error) {
	resp, err := s.GetResponsePaginate(AccountAPI, "organizations", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var data ScalewayOrganizationsDefinition

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// GetUser returns the user
func (s *ScalewayAPI) GetUser() (*ScalewayUserDefinition, error) {
	userID, err := s.GetUserID()
	if err != nil {
		return nil, err
	}
	resp, err := s.GetResponsePaginate(AccountAPI, fmt.Sprintf("users/%s", userID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var user ScalewayUsersDefinition

	if err = json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user.User, nil
}

// GetPermissions returns the permissions
func (s *ScalewayAPI) GetPermissions() (*ScalewayPermissionDefinition, error) {
	resp, err := s.GetResponsePaginate(AccountAPI, fmt.Sprintf("tokens/%s/permissions", s.Token), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var permissions ScalewayPermissionDefinition

	if err = json.Unmarshal(body, &permissions); err != nil {
		return nil, err
	}
	return &permissions, nil
}

// GetDashboard returns the dashboard
func (s *ScalewayAPI) GetDashboard() (*ScalewayDashboard, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "dashboard", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var dashboard ScalewayDashboardResp

	if err = json.Unmarshal(body, &dashboard); err != nil {
		return nil, err
	}
	return &dashboard.Dashboard, nil
}

func showResolverResults(needle string, results ScalewayResolverResults) error {
	w := tabwriter.NewWriter(os.Stderr, 20, 1, 3, ' ', 0)
	defer w.Flush()
	sort.Sort(results)
	fmt.Fprintf(w, "  IMAGEID\tFROM\tNAME\tZONE\tARCH\n")
	for _, result := range results {
		if result.Arch == "" {
			result.Arch = "n/a"
		}
		fmt.Fprintf(w, "- %s\t%s\t%s\t%s\t%s\n", result.TruncIdentifier(), result.CodeName(), result.Name, result.Region, result.Arch)
	}
	return fmt.Errorf("Too many candidates for %s (%d)", needle, len(results))
}

// GetContainers returns a ScalewayGetContainers
func (s *ScalewayAPI) GetContainers() (*ScalewayGetContainers, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "containers", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var containers ScalewayGetContainers

	if err = json.Unmarshal(body, &containers); err != nil {
		return nil, err
	}
	return &containers, nil
}

// GetContainerDatas returns a ScalewayGetContainerDatas
func (s *ScalewayAPI) GetContainerDatas(container string) (*ScalewayGetContainerDatas, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("containers/%s", container), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var datas ScalewayGetContainerDatas

	if err = json.Unmarshal(body, &datas); err != nil {
		return nil, err
	}
	return &datas, nil
}

// GetQuotas returns a ScalewayGetQuotas
func (s *ScalewayAPI) GetQuotas() (*ScalewayGetQuotas, error) {
	resp, err := s.GetResponsePaginate(AccountAPI, fmt.Sprintf("organizations/%s/quotas", s.Organization), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var quotas ScalewayGetQuotas

	if err = json.Unmarshal(body, &quotas); err != nil {
		return nil, err
	}
	return &quotas, nil
}

func rootNetDial(network, addr string) (net.Conn, error) {
	dialer := net.Dialer{
		Timeout:   10 * time.Second,
		KeepAlive: 10 * time.Second,
	}

	// bruteforce privileged ports
	var localAddr net.Addr
	var err error
	for port := 1; port <= 1024; port++ {
		localAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))

		// this should never happen
		if err != nil {
			return nil, err
		}

		dialer.LocalAddr = localAddr

		conn, err := dialer.Dial(network, addr)

		// if err is nil, dialer.Dial succeed, so let's go
		// else, err != nil, but we don't care
		if err == nil {
			return conn, nil
		}
	}
	// if here, all privileged ports were tried without success
	return nil, fmt.Errorf("bind: permission denied, are you root ?")
}

// SetPassword register the password
func (s *ScalewayAPI) SetPassword(password string) {
	s.password = password
}

// ResolveTTYUrl return an URL to get a tty
func (s *ScalewayAPI) ResolveTTYUrl() string {
	switch s.Region {
	case "par1", "":
		return "https://tty-par1.scaleway.com/v2/"
	case "ams1":
		return "https://tty-ams1.scaleway.com"
	}
	return ""
}
