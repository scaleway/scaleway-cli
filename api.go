package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// ScalewayAPI is the interface used to communicate with the Scaleway API
type ScalewayAPI struct {
	// APIEndpoint is the endpoint to the Scaleway API
	APIEndPoint string

	// Organization is the identifier of the Scaleway organization
	Organization string

	// Token is the authentication token for the Scaleway organization
	Token string

	// Cache is used to quickly resolve identifiers from names
	Cache *ScalewayCache
}

// ScalewayAPIError represents a Scaleway API Error
type ScalewayAPIError struct {
	// Message is a human-friendly error message
	APIMessage string `json:"message,omitempty"`

	// Type is a string code that defines the kind of error
	Type string `json:"type,omitempty"`

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
	Size int64 `json:"size,omitempty"`

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
	Public bool `json:"public",omitempty`

	// Bootscript is the bootscript bound to the image
	DefaultBootscript *ScalewayBootscript `json:"default_bootscript"`

	// Organization is the owner of the image
	Organization string `json:"organization"`

	// Arch is the architecture target of the image
	Arch string `json:"arch"`

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
	Size int64 `json:"size,omitempty"`

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
	Initrd ScalewayInitrd `json:initrd,omitempty`

	// Kernel is the kernel associated to this server
	Kernel ScalewayKernel `json:"kernel,omitempty"`

	// Public is true for public bootscripts and false for user bootscripts
	Public bool `json:"public",omitempty`
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

// ScalewaySecurityGroup represents a Scaleway security group
type ScalewaySecurityGroup struct {
	// Identifier is a unique identifier for the security group
	Identifier string `json:"id,omitempty"`

	// Name is the user-defined name of the security group
	Name string `json:"name,omitempty"`
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

	// PrivateIP reprensents the private IPV4 attached to the server (changes on each boot)
	PrivateIP string `json:"private_ip,omitempty"`

	// Bootscript is the unique identifier of the selected bootscript
	Bootscript *string `json:"bootscript"`

	// Hostname represents the ServerName in a format compatible with unix's hostname
	Hostname string `json:"hostname,omitempty"`

	// Tags represents user-defined tags
	Tags []string `json:"tags",omitempty`

	// Volumes are the attached volumes
	Volumes map[string]ScalewayVolume `json:"volumes,omitempty"`

	// SecurityGroup is the selected security group object
	SecurityGroup ScalewaySecurityGroup `json:"security_group,omitempty"`

	// Organization is the owner of the server
	Organization string `json:"organization"`
}

// ScalewayServerPathNameDefinition represents a Scaleway C1 server with only its name as field
type ScalewayServerPathNameDefinition struct {
	// Name is the user-defined name of the server
	Name string `json:"name"`
}

// ScalewayServerDefinition represents a Scaleway C1 server with image definition
type ScalewayServerDefinition struct {
	// Name is the user-defined name of the server
	Name string `json:"name"`

	// Image is the image used by the server
	Image *string `json:"image,omitempty"`

	// Volumes are the attached volumes
	Volumes map[string]string `json:"volumes,omitempty"`

	// Bootscript is the bootscript used by the server
	Bootscript *string `json:"bootscript"`

	// Tags are the metadata tags attached to the server
	Tags []string `json:"tags",omitempty`

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
	SnapshotIDentifier string `json:"root_volume"`
	Name               string `json:"name,omitempty"`
	Organization       string `json:"organization"`
	Arch               string `json:"arch"`
}

// NewScalewayAPI creates a ready-to-use ScalewayAPI client
func NewScalewayAPI(endpoint, organization, token string) (*ScalewayAPI, error) {
	cache, err := NewScalewayCache()
	if err != nil {
		return nil, err
	}
	return &ScalewayAPI{
		APIEndPoint:  endpoint,
		Organization: organization,
		Token:        token,
		Cache:        cache,
	}, nil
}

// Sync flushes out the cache to the disk
func (s *ScalewayAPI) Sync() {
	s.Cache.Save()
}

// GetResponse returns an http.Response object for the requested resource
func (s *ScalewayAPI) GetResponse(resource string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIEndPoint, "/"), resource)
	log.Debugf("GET %s", uri)
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

// PostResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PostResponse(resource string, data interface{}) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIEndPoint, "/"), resource)
	client := &http.Client{}
	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	log.Debugf("POST %s payload=%s", uri, strings.TrimSpace(fmt.Sprintf("%s", payload)))
	req, err := http.NewRequest("POST", uri, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

// PatchResponse returns an http.Response object for the updated resource
func (s *ScalewayAPI) PatchResponse(resource string, data interface{}) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIEndPoint, "/"), resource)
	client := &http.Client{}
	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	log.Debugf("PATCH %s payload=%s", uri, strings.TrimSpace(fmt.Sprintf("%s", payload)))
	req, err := http.NewRequest("PATCH", uri, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

// DeleteResponse returns an http.Response object for the deleted resource
func (s *ScalewayAPI) DeleteResponse(resource string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", strings.TrimRight(s.APIEndPoint, "/"), resource)
	client := &http.Client{}
	log.Debugf("DELETE %s", uri)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
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
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&servers)
	if err != nil {
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
	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != 200 {
		var error ScalewayAPIError
		err = decoder.Decode(&error)
		if err != nil {
			return nil, err
		}
		error.Debug()
		return nil, error
	}

	var oneServer ScalewayOneServer
	err = decoder.Decode(&oneServer)
	if err != nil {
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

	// Succeed POST code
	if resp.StatusCode == 202 {
		return nil
	}

	var error ScalewayAPIError
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&error)
	if err != nil {
		return err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return error
}

// DeleteServer deletes a server
func (s *ScalewayAPI) DeleteServer(serverID string) error {
	resp, err := s.DeleteResponse(fmt.Sprintf("servers/%s", serverID))
	if err != nil {
		return err
	}

	// Succeed POST code
	if resp.StatusCode == 204 {
		s.Cache.RemoveServer(serverID)
		return nil
	}

	var error ScalewayAPIError
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&error)
	if err != nil {
		return err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return error
}

// PostServer creates a new server
func (s *ScalewayAPI) PostServer(definition ScalewayServerDefinition) (string, error) {
	definition.Organization = s.Organization

	resp, err := s.PostResponse(fmt.Sprintf("servers"), definition)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	// Succeed POST code
	if resp.StatusCode == 201 {
		var server ScalewayOneServer
		err = decoder.Decode(&server)
		if err != nil {
			return "", err
		}
		s.Cache.InsertServer(server.Server.Identifier, server.Server.Name)
		return server.Server.Identifier, nil
	}

	var error ScalewayAPIError
	err = decoder.Decode(&error)

	if err != nil {
		return "", err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return "", error
}

// PatchServerName changes the name of the server
func (s *ScalewayAPI) PatchServerName(serverID string, definition ScalewayServerPathNameDefinition) error {
	resp, err := s.PatchResponse(fmt.Sprintf("servers/%s", serverID), definition)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	// Succeed PATCH code
	if resp.StatusCode == 200 {
		return nil
	}

	var error ScalewayAPIError
	err = decoder.Decode(&error)
	if err != nil {
		return err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return error
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
	decoder := json.NewDecoder(resp.Body)

	// Succeed POST code
	if resp.StatusCode == 201 {
		var snapshot ScalewayOneSnapshot
		err = decoder.Decode(&snapshot)
		if err != nil {
			return "", err
		}
		s.Cache.InsertSnapshot(snapshot.Snapshot.Identifier, snapshot.Snapshot.Name)
		return snapshot.Snapshot.Identifier, nil
	}

	var error ScalewayAPIError
	err = decoder.Decode(&error)

	if err != nil {
		return "", err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return "", error
}

// PostImage creates a new image
func (s *ScalewayAPI) PostImage(volumeID string, name string) (string, error) {
	definition := ScalewayImageDefinition{
		SnapshotIDentifier: volumeID,
		Name:               name,
		Organization:       s.Organization,
		Arch:               "arm",
	}
	resp, err := s.PostResponse("images", definition)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	// Succeed POST code
	if resp.StatusCode == 201 {
		var image ScalewayOneImage
		err = decoder.Decode(&image)
		if err != nil {
			return "", err
		}
		s.Cache.InsertImage(image.Image.Identifier, image.Image.Name)
		return image.Image.Identifier, nil
	}

	var error ScalewayAPIError
	err = decoder.Decode(&error)

	if err != nil {
		return "", err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return "", error
}

// PostVolume creates a new volume
func (s *ScalewayAPI) PostVolume(definition ScalewayVolumeDefinition) (string, error) {
	definition.Organization = s.Organization
	if definition.Type == "" {
		definition.Type = "l_ssd"
	}
	resp, err := s.PostResponse(fmt.Sprintf("volumes"), definition)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	// Succeed POST code
	if resp.StatusCode == 201 {
		var volume ScalewayOneVolume
		err = decoder.Decode(&volume)
		if err != nil {
			return "", err
		}
		// FIXME: s.Cache.InsertVolume(volume.Volume.Identifier, volume.Volume.Name)
		return volume.Volume.Identifier, nil
	}

	var error ScalewayAPIError
	err = decoder.Decode(&error)

	if err != nil {
		return "", err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return "", error
}

// ResolveServer attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveServer(needle string) ([]string, error) {
	servers := s.Cache.LookUpServers(needle)
	if len(servers) == 0 {
		_, err := s.GetServers(true, 0)
		if err != nil {
			return nil, err
		}
		servers = s.Cache.LookUpServers(needle)
	}
	return servers, nil
}

// ResolveSnapshot attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveSnapshot(needle string) ([]string, error) {
	snapshots := s.Cache.LookUpSnapshots(needle)
	if len(snapshots) == 0 {
		_, err := s.GetSnapshots()
		if err != nil {
			return nil, err
		}
		snapshots = s.Cache.LookUpSnapshots(needle)
	}
	return snapshots, nil
}

// ResolveImage attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveImage(needle string) ([]string, error) {
	images := s.Cache.LookUpImages(needle)
	if len(images) == 0 {
		_, err := s.GetImages()
		if err != nil {
			return nil, err
		}
		images = s.Cache.LookUpImages(needle)
	}
	return images, nil
}

// ResolveBootscript attempts the find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveBootscript(needle string) ([]string, error) {
	bootscripts := s.Cache.LookUpBootscripts(needle)
	if len(bootscripts) == 0 {
		_, err := s.GetBootscripts()
		if err != nil {
			return nil, err
		}
		bootscripts = s.Cache.LookUpBootscripts(needle)
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
	var images ScalewayImages
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&images)
	if err != nil {
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
	var oneImage ScalewayOneImage
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&oneImage)
	if err != nil {
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

	// Succeed POST code
	if resp.StatusCode == 204 {
		s.Cache.RemoveImage(imageID)
		return nil
	}

	var error ScalewayAPIError
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&error)
	if err != nil {
		return err
	}

	error.StatusCode = resp.StatusCode
	error.Debug()
	return error
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
	var snapshots ScalewaySnapshots
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&snapshots)
	if err != nil {
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
	var oneSnapshot ScalewayOneSnapshot
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&oneSnapshot)
	if err != nil {
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
	var volumes ScalewayVolumes
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&volumes)
	if err != nil {
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
	var oneVolume ScalewayOneVolume
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&oneVolume)
	if err != nil {
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
	var bootscripts ScalewayBootscripts
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&bootscripts)
	if err != nil {
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
	var oneBootscript ScalewayOneBootscript
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&oneBootscript)
	if err != nil {
		return nil, err
	}
	s.Cache.InsertBootscript(oneBootscript.Bootscript.Identifier, oneBootscript.Bootscript.Title)
	return &oneBootscript.Bootscript, nil
}

// GetTasks get the list of tasks from the ScalewayAPI
func (s *ScalewayAPI) GetTasks() (*[]ScalewayTask, error) {
	query := url.Values{}
	resp, err := s.GetResponse("tasks?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var tasks ScalewayTasks
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return &tasks.Tasks, nil
}

// CheckCredentials performs a dummy check to ensure we can contact the API
func (s *ScalewayAPI) CheckCredentials() error {
	query := url.Values{}
	query.Set("token_id", s.Token)
	resp, err := s.GetResponse("tokens?" + query.Encode())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Invalid credentials")
	}
	return nil
}
