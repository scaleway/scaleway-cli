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
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/vendor/github.com/moul/anonuuid"
	"github.com/scaleway/scaleway-cli/vendor/github.com/moul/http2curl"
)

// Default values
var (
	ComputeAPI  = "https://api.scaleway.com/"
	AccountAPI  = "https://account.scaleway.com/"
	MetadataAPI = "http://169.254.42.42/"
)

// ScalewayAPI is the interface used to communicate with the Scaleway API
type ScalewayAPI struct {
	// ComputeAPI is the endpoint to the Scaleway API
	ComputeAPI string

	// AccountAPI is the endpoint to the Scaleway Account API
	AccountAPI string

	// APIEndPoint or ACCOUNTEndPoint
	APIUrl string

	// Organization is the identifier of the Scaleway organization
	Organization string

	// Token is the authentication token for the Scaleway organization
	Token string

	// Password is the authentication password
	password string

	// Cache is used to quickly resolve identifiers from names
	Cache *ScalewayCache

	client *http.Client
	// Used when switching from an API to another
	oldTransport *http.RoundTripper
	anonuuid     anonuuid.AnonUUID
	isMetadata   bool
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
	if e.Message != "" {
		return e.Message
	}
	if e.APIMessage != "" {
		return e.APIMessage
	}
	if e.StatusCode != 0 {
		return fmt.Sprintf("Invalid return code, got %d", e.StatusCode)
	}
	panic(e)
}

// Debug create a debug log entry with HTTP error informations
func (e ScalewayAPIError) Debug() {
	log.WithFields(log.Fields{
		"StatusCode": e.StatusCode,
		"Type":       e.Type,
		"Message":    e.Message,
	}).Debug(e.APIMessage)

	// error.Fields handling
	for k, v := range e.Fields {
		log.Debugf("  %-30s %s", fmt.Sprintf("%s: ", k), v)
	}
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

// ScalewayVolumeDefinition represents a Scaleway C1 volume definition
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

// ScalewayVolumePutDefinition represents a Scaleway C1 volume with nullable fields (for PUT)
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

// ScalewayBootCmdArgs represents the boot arguments of a bootscript
type ScalewayBootCmdArgs struct {
	// Identifier is the unique identifier of boot args
	Identifier string `json:"id,omitempty"`

	// Value is the content of the cmd args
	Value string `json:"value,omitempty"`
}

// ScalewayInitrd represents the initrd used by a bootscript
type ScalewayInitrd struct {
	// Identifier is the unique identifier of the initrd
	Identifier string `json:"id,omitempty"`

	// Path is the path to the initrd used
	Path string `json:"path,omitempty"`

	// Title is the title of the initrd used
	Title string `json:"title,omitempty"`
}

// ScalewayKernel represents a kernel used on C1 servers
type ScalewayKernel struct {
	// Identifier is the unique identifier of the kernel
	Identifier string `json:"id,omitempty"`

	// DTB is the kernel DTB used by this kernel
	DTB string `json:"dtb"`

	// Path is the path to the kernel image
	Path string `json:"path,omitempty"`

	// Title is the title of the kernel
	Title string `json:"title,omitempty"`
}

// ScalewayBootscript represents a Scaleway Bootscript
type ScalewayBootscript struct {
	// Identifier is a unique identifier for the bootscript
	Identifier string `json:"id,omitempty"`

	// Name is a user-defined name for the bootscript
	Title string `json:"title,omitempty"`

	// BootCmdArgs represents the arguments used to boot
	BootCmdArgs ScalewayBootCmdArgs `json:"bootcmdargs,omitempty"`

	// Initrd is the initrd used by this bootscript
	Initrd ScalewayInitrd `json:"initrd,omitempty"`

	// Kernel is the kernel associated to this server
	Kernel ScalewayKernel `json:"kernel,omitempty"`

	// Public is true for public bootscripts and false for user bootscripts
	Public bool `json:"public,omitempty"`
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
	Postion      int    `json:"position"`
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
	EnableDefaultSecurity bool                    `json:"enable_default_security"`
	ID                    string                  `json:"id"`
	Organization          string                  `json:"organization"`
	Name                  string                  `json:"name"`
	OrganizationDefault   bool                    `json:"organization_default"`
	Servers               []ScalewaySecurityGroup `json:"servers"`
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
	Organization string `json:"organization"`
	Reverse      string `json:"reverse"`
	ID           string `json:"id"`
	Server       struct {
		Identifier string `json:"id,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"server,omitempty"`
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

// ScalewayNewSecurityGroup definition POST/PUT request /security_groups
type ScalewayNewSecurityGroup struct {
	Organization string `json:"organization"`
	Name         string `json:"name"`
	Description  string `json:"description"`
}

// ScalewayServer represents a Scaleway C1 server
type ScalewayServer struct {
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
}

// ScalewayServerPatchDefinition represents a Scaleway C1 server with nullable fields (for PATCH)
type ScalewayServerPatchDefinition struct {
	Name              *string                    `json:"name,omitempty"`
	CreationDate      *string                    `json:"creation_date,omitempty"`
	ModificationDate  *string                    `json:"modification_date,omitempty"`
	Image             *ScalewayImage             `json:"image,omitempty"`
	DynamicIPRequired *bool                      `json:"dynamic_ip_required,omitempty"`
	PublicAddress     *ScalewayIPAddress         `json:"public_ip,omitempty"`
	State             *string                    `json:"state,omitempty"`
	StateDetail       *string                    `json:"state_detail,omitempty"`
	PrivateIP         *string                    `json:"private_ip,omitempty"`
	Bootscript        *ScalewayBootscript        `json:"bootscript,omitempty"`
	Hostname          *string                    `json:"hostname,omitempty"`
	Volumes           *map[string]ScalewayVolume `json:"volumes,omitempty"`
	SecurityGroup     *ScalewaySecurityGroup     `json:"security_group,omitempty"`
	Organization      *string                    `json:"organization,omitempty"`
	Tags              *[]string                  `json:"tags,omitempty"`
}

// ScalewayServerDefinition represents a Scaleway C1 server with image definition
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
}

// ScalewayOneServer represents the response of a GET /servers/UUID API call
type ScalewayOneServer struct {
	Server ScalewayServer `json:"server,omitempty"`
}

// ScalewayServers represents a group of Scaleway C1 servers
type ScalewayServers struct {
	// Servers holds scaleway servers of the response
	Servers []ScalewayServer `json:"servers,omitempty"`
}

// ScalewayServerAction represents an action to perform on a Scaleway C1 server
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
	Key string `json:"key"`
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

// ScalewayUserdata represents []byte
type ScalewayUserdata []byte

// FuncMap used for json inspection
var FuncMap = template.FuncMap{
	"json": func(v interface{}) string {
		a, _ := json.Marshal(v)
		return string(a)
	},
}

// NewScalewayAPI creates a ready-to-use ScalewayAPI client
func NewScalewayAPI(apiEndPoint, accountEndPoint, organization, token string) (*ScalewayAPI, error) {
	cache, err := NewScalewayCache()
	if err != nil {
		return nil, err
	}
	s := &ScalewayAPI{
		// exposed
		ComputeAPI:   apiEndPoint,
		AccountAPI:   accountEndPoint,
		APIUrl:       apiEndPoint,
		Organization: organization,
		Token:        token,
		Cache:        cache,
		password:     "",

		// internal
		anonuuid: *anonuuid.New(),
		client:   &http.Client{},
	}

	if os.Getenv("SCALEWAY_TLSVERIFY") == "0" {
		s.client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return s, nil
}

// Sync flushes out the cache to the disk
func (s *ScalewayAPI) Sync() {
	s.Cache.Save()
}

// GetResponse returns an http.Response object for the requested resource
func (s *ScalewayAPI) GetResponse(resource string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIUrl, "/"), resource)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("SCW_SENSITIVE") != "1" {
		log.Debug(s.HideAPICredentials(curl.String()))
	} else {
		log.Debug(curl.String())
	}
	return s.client.Do(req)
}

// PostResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PostResponse(resource string, data interface{}) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIUrl, "/"), resource)
	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("SCW_SENSITIVE") != "1" {
		log.Debug(s.HideAPICredentials(curl.String()))
	} else {
		log.Debug(curl.String())
	}

	return s.client.Do(req)
}

// PatchResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PatchResponse(resource string, data interface{}) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIUrl, "/"), resource)
	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", uri, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("SCW_SENSITIVE") != "1" {
		log.Debug(s.HideAPICredentials(curl.String()))
	} else {
		log.Debug(curl.String())
	}

	return s.client.Do(req)
}

// PutResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PutResponse(resource string, data interface{}) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIUrl, "/"), resource)
	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", uri, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("SCW_SENSITIVE") != "1" {
		log.Debug(s.HideAPICredentials(curl.String()))
	} else {
		log.Debug(curl.String())
	}

	return s.client.Do(req)
}

// DeleteResponse returns an http.Response object for the deleted resource
func (s *ScalewayAPI) DeleteResponse(resource string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIUrl, "/"), resource)

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return nil, err
	}
	if os.Getenv("SCW_SENSITIVE") != "1" {
		log.Debug(s.HideAPICredentials(curl.String()))
	} else {
		log.Debug(curl.String())
	}

	return s.client.Do(req)
}

// handleHTTPError checks the statusCode and displays the error
func handleHTTPError(goodStatusCode []int, resp *http.Response) (*json.Decoder, error) {
	if resp.StatusCode >= 500 {
		errorMsg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(errorMsg))
	}
	good := false
	for _, code := range goodStatusCode {
		if code == resp.StatusCode {
			good = true
		}
	}
	decoder := json.NewDecoder(resp.Body)

	if !good {
		var scwError ScalewayAPIError

		err := decoder.Decode(&scwError)
		if err != nil {
			return nil, err
		}
		scwError.StatusCode = resp.StatusCode
		scwError.Debug()
		return nil, scwError
	}
	return decoder, nil
}

// GetServers gets the list of servers from the ScalewayAPI
func (s *ScalewayAPI) GetServers(all bool, limit int) (*[]ScalewayServer, error) {
	query := url.Values{}
	if !all {
		query.Set("state", "running")
	}
	if limit > 0 {
		// FIXME: wait for the API to be ready
		// query.Set("per_page", strconv.Itoa(limit))
	}
	if all && limit == 0 {
		s.Cache.ClearServers()
	}
	resp, err := s.GetResponse("servers?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var servers ScalewayServers

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	if err = decoder.Decode(&servers); err != nil {
		return nil, err
	}
	for _, server := range servers.Servers {
		s.Cache.InsertServer(server.Identifier, server.Name)
	}
	// FIXME: when API limit is ready, remove the following code
	if limit > 0 && limit < len(servers.Servers) {
		servers.Servers = servers.Servers[0:limit]
	}
	return &servers.Servers, nil
}

// GetServer gets a server from the ScalewayAPI
func (s *ScalewayAPI) GetServer(serverID string) (*ScalewayServer, error) {
	resp, err := s.GetResponse("servers/" + serverID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}

	var oneServer ScalewayOneServer

	if err = decoder.Decode(&oneServer); err != nil {
		return nil, err
	}
	s.Cache.InsertServer(oneServer.Server.Identifier, oneServer.Server.Name)
	return &oneServer.Server, nil
}

// PostServerAction posts an action on a server
func (s *ScalewayAPI) PostServerAction(serverID, action string) error {
	data := ScalewayServerAction{
		Action: action,
	}
	resp, err := s.PostResponse(fmt.Sprintf("servers/%s/action", serverID), data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = handleHTTPError([]int{202}, resp); err != nil {
		return err
	}
	return nil
}

// DeleteServer deletes a server
func (s *ScalewayAPI) DeleteServer(serverID string) error {
	resp, err := s.DeleteResponse(fmt.Sprintf("servers/%s", serverID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = handleHTTPError([]int{204}, resp); err != nil {
		return err
	}
	s.Cache.RemoveServer(serverID)
	return nil
}

// PostServer creates a new server
func (s *ScalewayAPI) PostServer(definition ScalewayServerDefinition) (string, error) {
	definition.Organization = s.Organization

	resp, err := s.PostResponse("servers", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{201}, resp)
	if err != nil {
		return "", err
	}
	var server ScalewayOneServer

	if err = decoder.Decode(&server); err != nil {
		return "", err
	}
	s.Cache.InsertServer(server.Server.Identifier, server.Server.Name)
	return server.Server.Identifier, nil
}

// PatchUserSSHKey updates a user
func (s *ScalewayAPI) PatchUserSSHKey(UserID string, definition ScalewayUserPatchSSHKeyDefinition) error {
	s.EnableAccountAPI()
	defer s.DisableAccountAPI()
	resp, err := s.PatchResponse(fmt.Sprintf("users/%s", UserID), definition)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := handleHTTPError([]int{200}, resp); err != nil {
		return err
	}
	return nil
}

// PatchServer updates a server
func (s *ScalewayAPI) PatchServer(serverID string, definition ScalewayServerPatchDefinition) error {
	resp, err := s.PatchResponse(fmt.Sprintf("servers/%s", serverID), definition)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := handleHTTPError([]int{200}, resp); err != nil {
		return err
	}
	return nil
}

// PostSnapshot creates a new snapshot
func (s *ScalewayAPI) PostSnapshot(volumeID string, name string) (string, error) {
	definition := ScalewaySnapshotDefinition{
		VolumeIDentifier: volumeID,
		Name:             name,
		Organization:     s.Organization,
	}
	resp, err := s.PostResponse("snapshots", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{201}, resp)
	if err != nil {
		return "", err
	}
	var snapshot ScalewayOneSnapshot

	if err = decoder.Decode(&snapshot); err != nil {
		return "", err
	}
	s.Cache.InsertSnapshot(snapshot.Snapshot.Identifier, snapshot.Snapshot.Name)
	return snapshot.Snapshot.Identifier, nil
}

// PostImage creates a new image
func (s *ScalewayAPI) PostImage(volumeID string, name string, bootscript string) (string, error) {
	definition := ScalewayImageDefinition{
		SnapshotIDentifier: volumeID,
		Name:               name,
		Organization:       s.Organization,
		Arch:               "arm",
	}
	if bootscript != "" {
		definition.DefaultBootscript = &bootscript
	}

	resp, err := s.PostResponse("images", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{201}, resp)
	if err != nil {
		return "", err
	}
	var image ScalewayOneImage

	if err = decoder.Decode(&image); err != nil {
		return "", err
	}
	s.Cache.InsertImage(image.Image.Identifier, image.Image.Name)
	return image.Image.Identifier, nil
}

// PostVolume creates a new volume
func (s *ScalewayAPI) PostVolume(definition ScalewayVolumeDefinition) (string, error) {
	definition.Organization = s.Organization
	if definition.Type == "" {
		definition.Type = "l_ssd"
	}

	resp, err := s.PostResponse("volumes", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{201}, resp)
	if err != nil {
		return "", err
	}
	var volume ScalewayOneVolume

	if err = decoder.Decode(&volume); err != nil {
		return "", err
	}
	// FIXME: s.Cache.InsertVolume(volume.Volume.Identifier, volume.Volume.Name)
	return volume.Volume.Identifier, nil
}

// PutVolume updates a volume
func (s *ScalewayAPI) PutVolume(volumeID string, definition ScalewayVolumePutDefinition) error {
	resp, err := s.PutResponse(fmt.Sprintf("volumes/%s", volumeID), definition)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := handleHTTPError([]int{200}, resp); err != nil {
		return err
	}
	return nil
}

// ResolveServer attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveServer(needle string) (ScalewayResolverResults, error) {
	servers := s.Cache.LookUpServers(needle, true)
	if len(servers) == 0 {
		if _, err := s.GetServers(true, 0); err != nil {
			return nil, err
		}
		servers = s.Cache.LookUpServers(needle, true)
	}
	return servers, nil
}

// ResolveSnapshot attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveSnapshot(needle string) (ScalewayResolverResults, error) {
	snapshots := s.Cache.LookUpSnapshots(needle, true)
	if len(snapshots) == 0 {
		if _, err := s.GetSnapshots(); err != nil {
			return nil, err
		}
		snapshots = s.Cache.LookUpSnapshots(needle, true)
	}
	return snapshots, nil
}

// ResolveImage attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveImage(needle string) (ScalewayResolverResults, error) {
	images := s.Cache.LookUpImages(needle, true)
	if len(images) == 0 {
		if _, err := s.GetImages(); err != nil {
			return nil, err
		}
		images = s.Cache.LookUpImages(needle, true)
	}
	return images, nil
}

// ResolveBootscript attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveBootscript(needle string) (ScalewayResolverResults, error) {
	bootscripts := s.Cache.LookUpBootscripts(needle, true)
	if len(bootscripts) == 0 {
		if _, err := s.GetBootscripts(); err != nil {
			return nil, err
		}
		bootscripts = s.Cache.LookUpBootscripts(needle, true)
	}
	return bootscripts, nil
}

// GetImages gets the list of images from the ScalewayAPI
func (s *ScalewayAPI) GetImages() (*[]ScalewayImage, error) {
	query := url.Values{}
	s.Cache.ClearImages()
	resp, err := s.GetResponse("images?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var images ScalewayImages

	if err = decoder.Decode(&images); err != nil {
		return nil, err
	}
	for _, image := range images.Images {
		s.Cache.InsertImage(image.Identifier, image.Name)
	}
	return &images.Images, nil
}

// GetImage gets an image from the ScalewayAPI
func (s *ScalewayAPI) GetImage(imageID string) (*ScalewayImage, error) {
	resp, err := s.GetResponse("images/" + imageID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var oneImage ScalewayOneImage

	if err = decoder.Decode(&oneImage); err != nil {
		return nil, err
	}
	s.Cache.InsertImage(oneImage.Image.Identifier, oneImage.Image.Name)
	return &oneImage.Image, nil
}

// DeleteImage deletes a image
func (s *ScalewayAPI) DeleteImage(imageID string) error {
	resp, err := s.DeleteResponse(fmt.Sprintf("images/%s", imageID))
	if err != nil {
		s.Cache.RemoveImage(imageID)
		return err
	}
	defer resp.Body.Close()

	if _, err := handleHTTPError([]int{204}, resp); err != nil {
		return err
	}
	s.Cache.RemoveImage(imageID)
	return nil
}

// GetSnapshots gets the list of snapshots from the ScalewayAPI
func (s *ScalewayAPI) GetSnapshots() (*[]ScalewaySnapshot, error) {
	query := url.Values{}
	s.Cache.ClearSnapshots()

	resp, err := s.GetResponse("snapshots?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var snapshots ScalewaySnapshots

	if err = decoder.Decode(&snapshots); err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots.Snapshots {
		s.Cache.InsertSnapshot(snapshot.Identifier, snapshot.Name)
	}
	return &snapshots.Snapshots, nil
}

// GetSnapshot gets a snapshot from the ScalewayAPI
func (s *ScalewayAPI) GetSnapshot(snapshotID string) (*ScalewaySnapshot, error) {
	resp, err := s.GetResponse("snapshots/" + snapshotID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var oneSnapshot ScalewayOneSnapshot

	if err = decoder.Decode(&oneSnapshot); err != nil {
		return nil, err
	}
	s.Cache.InsertSnapshot(oneSnapshot.Snapshot.Identifier, oneSnapshot.Snapshot.Name)
	return &oneSnapshot.Snapshot, nil
}

// GetVolumes gets the list of volumes from the ScalewayAPI
func (s *ScalewayAPI) GetVolumes() (*[]ScalewayVolume, error) {
	query := url.Values{}
	s.Cache.ClearVolumes()

	resp, err := s.GetResponse("volumes?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var volumes ScalewayVolumes

	if err = decoder.Decode(&volumes); err != nil {
		return nil, err
	}
	for _, volume := range volumes.Volumes {
		s.Cache.InsertVolume(volume.Identifier, volume.Name)
	}
	return &volumes.Volumes, nil
}

// GetVolume gets a volume from the ScalewayAPI
func (s *ScalewayAPI) GetVolume(volumeID string) (*ScalewayVolume, error) {
	resp, err := s.GetResponse("volumes/" + volumeID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var oneVolume ScalewayOneVolume

	if err = decoder.Decode(&oneVolume); err != nil {
		return nil, err
	}
	s.Cache.InsertVolume(oneVolume.Volume.Identifier, oneVolume.Volume.Name)
	return &oneVolume.Volume, nil
}

// GetBootscripts gets the list of bootscripts from the ScalewayAPI
func (s *ScalewayAPI) GetBootscripts() (*[]ScalewayBootscript, error) {
	query := url.Values{}
	s.Cache.ClearBootscripts()
	resp, err := s.GetResponse("bootscripts?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var bootscripts ScalewayBootscripts

	if err = decoder.Decode(&bootscripts); err != nil {
		return nil, err
	}
	for _, bootscript := range bootscripts.Bootscripts {
		s.Cache.InsertBootscript(bootscript.Identifier, bootscript.Title)
	}
	return &bootscripts.Bootscripts, nil
}

// GetBootscript gets a bootscript from the ScalewayAPI
func (s *ScalewayAPI) GetBootscript(bootscriptID string) (*ScalewayBootscript, error) {
	resp, err := s.GetResponse("bootscripts/" + bootscriptID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var oneBootscript ScalewayOneBootscript

	if err = decoder.Decode(&oneBootscript); err != nil {
		return nil, err
	}
	s.Cache.InsertBootscript(oneBootscript.Bootscript.Identifier, oneBootscript.Bootscript.Title)
	return &oneBootscript.Bootscript, nil
}

// GetUserdatas gets list of userdata for a server
func (s *ScalewayAPI) GetUserdatas(serverID string) (*ScalewayUserdatas, error) {
	var url string
	if s.isMetadata {
		url = "/user_data"
	} else {
		url = fmt.Sprintf("servers/%s/user_data", serverID)
	}

	resp, err := s.GetResponse(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var userdatas ScalewayUserdatas

	if err = decoder.Decode(&userdatas); err != nil {
		return nil, err
	}
	return &userdatas, nil
}

func (s *ScalewayUserdata) String() string {
	return string(*s)
}

// GetUserdata gets a specific userdata for a server
func (s *ScalewayAPI) GetUserdata(serverID string, key string) (*ScalewayUserdata, error) {
	var data ScalewayUserdata
	var err error
	var url string

	if s.isMetadata {
		url = fmt.Sprintf("/user_data/%s", key)
	} else {
		url = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	resp, err := s.GetResponse(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("no such user_data %q (%d)", key, resp.StatusCode)
	}
	data, err = ioutil.ReadAll(resp.Body)
	return &data, err
}

// PatchUserdata sets a user data
func (s *ScalewayAPI) PatchUserdata(serverID string, key string, value []byte) error {
	var resource string

	if s.isMetadata {
		resource = fmt.Sprintf("/user_data/%s", key)
	} else {
		resource = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIUrl, "/"), resource)
	payload := new(bytes.Buffer)
	payload.Write(value)

	req, err := http.NewRequest("PATCH", uri, payload)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "text/plain")

	curl, err := http2curl.GetCurlCommand(req)
	if os.Getenv("SCW_SENSITIVE") != "1" {
		log.Debug(s.HideAPICredentials(curl.String()))
	} else {
		log.Debug(curl.String())
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return nil
	}

	return fmt.Errorf("cannot set user_data (%d)", resp.StatusCode)
}

// DeleteUserdata deletes a server user_data
func (s *ScalewayAPI) DeleteUserdata(serverID string, key string) error {
	var url string

	if s.isMetadata {
		url = fmt.Sprintf("/user_data/%s", key)
	} else {
		url = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	resp, err := s.DeleteResponse(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Succeed POST code
	if resp.StatusCode == 204 {
		return nil
	}
	return fmt.Errorf("cannot delete user_data (%d)", resp.StatusCode)
}

// GetTasks get the list of tasks from the ScalewayAPI
func (s *ScalewayAPI) GetTasks() (*[]ScalewayTask, error) {
	query := url.Values{}
	resp, err := s.GetResponse("tasks?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var tasks ScalewayTasks

	if err = decoder.Decode(&tasks); err != nil {
		return nil, err
	}
	return &tasks.Tasks, nil
}

// CheckCredentials performs a dummy check to ensure we can contact the API
func (s *ScalewayAPI) CheckCredentials() error {
	s.EnableAccountAPI()
	defer s.DisableAccountAPI()
	query := url.Values{}
	query.Set("token_id", s.Token)

	resp, err := s.GetResponse("tokens?" + query.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := handleHTTPError([]int{200}, resp); err != nil {
		return err
	}
	return nil
}

// GetUserID returns the userID
func (s *ScalewayAPI) GetUserID() (string, error) {
	s.EnableAccountAPI()
	defer s.DisableAccountAPI()

	resp, err := s.GetResponse(fmt.Sprintf("tokens/%s", s.Token))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return "", err
	}
	var token ScalewayTokensDefinition

	if err = decoder.Decode(&token); err != nil {
		return "", err
	}
	return token.Token.UserID, nil
}

// GetOrganization returns Organization
func (s *ScalewayAPI) GetOrganization() (*ScalewayOrganizationsDefinition, error) {
	s.EnableAccountAPI()
	defer s.DisableAccountAPI()

	resp, err := s.GetResponse("organizations")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var data ScalewayOrganizationsDefinition

	if err = decoder.Decode(&data); err != nil {
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
	s.EnableAccountAPI()
	defer s.DisableAccountAPI()

	resp, err := s.GetResponse(fmt.Sprintf("users/%s", userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var user ScalewayUsersDefinition

	if err = decoder.Decode(&user); err != nil {
		return nil, err
	}
	return &user.User, nil
}

// GetPermissions returns the permissions
func (s *ScalewayAPI) GetPermissions() (*ScalewayPermissionDefinition, error) {
	s.EnableAccountAPI()
	defer s.DisableAccountAPI()
	resp, err := s.GetResponse(fmt.Sprintf("tokens/%s/permissions", s.Token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var permissions ScalewayPermissionDefinition

	if err = decoder.Decode(&permissions); err != nil {
		return nil, err
	}
	return &permissions, nil
}

// GetDashboard returns the dashboard
func (s *ScalewayAPI) GetDashboard() (*ScalewayDashboard, error) {
	resp, err := s.GetResponse("dashboard")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var dashboard ScalewayDashboardResp

	if err = decoder.Decode(&dashboard); err != nil {
		return nil, err
	}
	return &dashboard.Dashboard, nil
}

// GetServerID returns exactly one server matching or dies
func (s *ScalewayAPI) GetServerID(needle string) string {
	// Parses optional type prefix, i.e: "server:name" -> "name"
	_, needle = parseNeedle(needle)

	servers, err := s.ResolveServer(needle)
	if err != nil {
		log.Fatalf("Unable to resolve server %s: %s", needle, err)
	}
	if len(servers) == 1 {
		return servers[0].Identifier
	}
	if len(servers) == 0 {
		log.Fatalf("No such server: %s", needle)
	}

	showResolverResults(needle, servers)
	os.Exit(1)
	return ""
}

func showResolverResults(needle string, results ScalewayResolverResults) error {
	log.Errorf("Too many candidates for %s (%d)", needle, len(results))

	w := tabwriter.NewWriter(os.Stderr, 20, 1, 3, ' ', 0)
	defer w.Flush()
	sort.Sort(results)
	for _, result := range results {
		fmt.Fprintf(w, "- %s\t%s\t%s\n", result.TruncIdentifier(), result.CodeName(), result.Name)
	}
	return nil
}

// GetSnapshotID returns exactly one snapshot matching or dies
func (s *ScalewayAPI) GetSnapshotID(needle string) string {
	// Parses optional type prefix, i.e: "snapshot:name" -> "name"
	_, needle = parseNeedle(needle)

	snapshots, err := s.ResolveSnapshot(needle)
	if err != nil {
		log.Fatalf("Unable to resolve snapshot %s: %s", needle, err)
	}
	if len(snapshots) == 1 {
		return snapshots[0].Identifier
	}
	if len(snapshots) == 0 {
		log.Fatalf("No such snapshot: %s", needle)
	}

	showResolverResults(needle, snapshots)
	os.Exit(1)
	return ""
}

// GetImageID returns exactly one image matching or dies
func (s *ScalewayAPI) GetImageID(needle string, exitIfMissing bool) string {
	// Parses optional type prefix, i.e: "image:name" -> "name"
	_, needle = parseNeedle(needle)

	images, err := s.ResolveImage(needle)
	if err != nil {
		log.Fatalf("Unable to resolve image %s: %s", needle, err)
	}
	if len(images) == 1 {
		return images[0].Identifier
	}
	if len(images) == 0 {
		if exitIfMissing {
			log.Fatalf("No such image: %s", needle)
		} else {
			return ""
		}
	}

	showResolverResults(needle, images)
	os.Exit(1)
	return ""
}

// GetSecurityGroups returns a ScalewaySecurityGroups
func (s *ScalewayAPI) GetSecurityGroups() (*ScalewayGetSecurityGroups, error) {
	resp, err := s.GetResponse("security_groups")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups ScalewayGetSecurityGroups

	if err = decoder.Decode(&securityGroups); err != nil {
		return nil, err
	}
	return &securityGroups, nil
}

// GetSecurityGroupRules returns a ScalewaySecurityGroupRules
func (s *ScalewayAPI) GetSecurityGroupRules(groupID string) (*ScalewayGetSecurityGroupRules, error) {
	resp, err := s.GetResponse(fmt.Sprintf("security_groups/%s/rules", groupID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroupRules ScalewayGetSecurityGroupRules

	if err = decoder.Decode(&securityGroupRules); err != nil {
		return nil, err
	}
	return &securityGroupRules, nil
}

// GetASecurityGroupRule returns a ScalewaySecurityGroupRule
func (s *ScalewayAPI) GetASecurityGroupRule(groupID string, rulesID string) (*ScalewayGetSecurityGroupRule, error) {
	resp, err := s.GetResponse(fmt.Sprintf("security_groups/%s/rules/%s", groupID, rulesID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroupRules ScalewayGetSecurityGroupRule

	if err = decoder.Decode(&securityGroupRules); err != nil {
		return nil, err
	}
	return &securityGroupRules, nil
}

// GetASecurityGroup returns a ScalewaySecurityGroup
func (s *ScalewayAPI) GetASecurityGroup(groupsID string) (*ScalewayGetSecurityGroup, error) {
	resp, err := s.GetResponse(fmt.Sprintf("security_groups/%s", groupsID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups ScalewayGetSecurityGroup

	if err = decoder.Decode(&securityGroups); err != nil {
		return nil, err
	}
	return &securityGroups, nil
}

// GetIPS returns a ScalewayGetIPS
func (s *ScalewayAPI) GetIPS() (*ScalewayGetIPS, error) {
	resp, err := s.GetResponse("ips")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{200}, resp)
	if err != nil {
		return nil, err
	}
	var ips ScalewayGetIPS

	if err = decoder.Decode(&ips); err != nil {
		return nil, err
	}
	return &ips, nil
}

// NewIP returns a new IP
func (s *ScalewayAPI) NewIP() (*ScalewayGetIP, error) {
	var orga struct {
		Organization string `json:"organization"`
	}
	orga.Organization = s.Organization
	resp, err := s.PostResponse("ips", orga)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{201}, resp)
	if err != nil {
		return nil, err
	}
	var ip ScalewayGetIP

	if err = decoder.Decode(&ip); err != nil {
		return nil, err
	}
	return &ip, nil
}

// AttachIP attachs an IP to a server
func (s *ScalewayAPI) AttachIP(ipID, serverID string) error {
	var update struct {
		Address      string  `json:"address"`
		ID           string  `json:"id"`
		Reverse      *string `json:"reverse"`
		Organization string  `json:"organization"`
		Server       string  `json:"server"`
	}

	ip, err := s.GetIP(ipID)
	if err != nil {
		return err
	}
	update.Address = ip.IP.Address
	update.ID = ip.IP.ID
	update.Organization = ip.IP.Organization
	update.Server = serverID
	resp, err := s.PutResponse(fmt.Sprintf("ips/%s", ipID), update)
	if err != nil {
		return err
	}
	if _, err := handleHTTPError([]int{200}, resp); err != nil {
		return err
	}
	return nil
}

// DeleteIP deletes an IP
func (s *ScalewayAPI) DeleteIP(ipID string) error {
	resp, err := s.DeleteResponse(fmt.Sprintf("ips/%s", ipID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := handleHTTPError([]int{204}, resp); err != nil {
		return err
	}
	return nil
}

// GetIP returns a ScalewayGetIP
func (s *ScalewayAPI) GetIP(ipID string) (*ScalewayGetIP, error) {
	resp, err := s.GetResponse(fmt.Sprintf("ips/%s", ipID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder, err := handleHTTPError([]int{204}, resp)
	if err != nil {
		return nil, err
	}
	var ip ScalewayGetIP

	if err = decoder.Decode(&ip); err != nil {
		return nil, err
	}
	return &ip, nil
}

// GetBootscriptID returns exactly one bootscript matching or dies
func (s *ScalewayAPI) GetBootscriptID(needle string) string {
	// Parses optional type prefix, i.e: "bootscript:name" -> "name"
	_, needle = parseNeedle(needle)

	bootscripts, err := s.ResolveBootscript(needle)
	if err != nil {
		log.Fatalf("Unable to resolve bootscript %s: %s", needle, err)
	}
	if len(bootscripts) == 1 {
		return bootscripts[0].Identifier
	}
	if len(bootscripts) == 0 {
		log.Fatalf("No such bootscript: %s", needle)
	}

	showResolverResults(needle, bootscripts)
	os.Exit(1)
	return ""
}

// HideAPICredentials removes API credentials from a string
func (s *ScalewayAPI) HideAPICredentials(input string) string {
	output := input
	if s.Token != "" {
		output = strings.Replace(output, s.Token, s.anonuuid.FakeUUID(s.Token), -1)
	}
	if s.Organization != "" {
		output = strings.Replace(output, s.Organization, s.anonuuid.FakeUUID(s.Organization), -1)
	}
	if s.password != "" {
		output = strings.Replace(output, s.password, "XX-XX-XX-XX", -1)
	}
	return output
}

// EnableAccountAPI enable accountAPI
func (s *ScalewayAPI) EnableAccountAPI() {
	s.APIUrl = s.AccountAPI
}

// DisableAccountAPI disable accountAPI
func (s *ScalewayAPI) DisableAccountAPI() {
	s.APIUrl = s.ComputeAPI
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

// EnableMetadataAPI enable metadataAPI
func (s *ScalewayAPI) EnableMetadataAPI() {
	s.APIUrl = MetadataAPI
	if os.Getenv("SCW_METADATA_URL") != "" {
		s.APIUrl = os.Getenv("SCW_METADATA_URL")
	}
	s.oldTransport = &s.client.Transport
	s.client.Transport = &http.Transport{
		Dial: rootNetDial,
	}
	s.isMetadata = true
}

// DisableMetadataAPI disable metadataAPI
func (s *ScalewayAPI) DisableMetadataAPI() {
	s.APIUrl = s.ComputeAPI
	s.client.Transport = *s.oldTransport
	s.isMetadata = false
}

// SetPassword register the password
func (s *ScalewayAPI) SetPassword(password string) {
	s.password = password
}
