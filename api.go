package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ScalewayAPI is the interface used to communicate with the Scaleway API
type ScalewayAPI struct {
	// APIEndpoint is the endpoint to the Scaleway API
	APIEndPoint string

	// Organization is the identifier of the Scaleway orgnization
	Organization string

	// Token is the authentication token for the Scaleway organization
	Token string
}

// ScalewayIPAddress represents a Scaleway IP address
type ScalewayIPAddress struct {
	// IP is an IPv4 address
	IP string `json:"address,omitempty"`
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

	// PublicIP is the public IP address bound to the server
	PublicAddress ScalewayIPAddress `json:"public_ip,omitempty"`

	// State is the current status of the server
	State string `json:"state,omitempty"`
}

// ShortIdentifier returns a truncated identifier
func (s *ScalewayServer) ShortIdentifier() string {
	if len(s.Identifier) > 8 {
		return s.Identifier[:8]
	}
	return s.Identifier
}

// ScalewayServers represents a group of Scaleway C1 servers
type ScalewayServers struct {
	// Servers holds scaleway servers of the response
	Servers []ScalewayServer `json:"servers,omitempty"`
}

// NewScalewayAPI creates a ready-to-use ScalewayAPI client
func NewScalewayAPI(endpoint, organization, token string) *ScalewayAPI {
	return &ScalewayAPI{
		APIEndPoint:  endpoint,
		Organization: organization,
		Token:        token,
	}
}

// GetResponse returns a http.Response object for the requested resource
func (s *ScalewayAPI) GetResponse(resource string) (*http.Response, error) {
	uri := fmt.Sprintf("%s/%s", s.APIEndPoint, resource)
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

// GetServers get the list of servers from the ScalewayAPI
func (s *ScalewayAPI) GetServers() (*[]ScalewayServer, error) {
	resp, err := s.GetResponse("servers")
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
	return &servers.Servers, nil
}
