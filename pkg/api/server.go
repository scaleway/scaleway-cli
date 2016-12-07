// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

// Interact with Scaleway API

// Package api contains client and functions to interact with Scaleway API
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/sync/errgroup"
)

type ServerAPI interface {
	GetServer(serverID string) (*ScalewayServer, error)
	GetServers(all bool, limit int) (*[]ScalewayServer, error)
	PostServer(definition ScalewayServerDefinition) (string, error)
	PostServerAction(serverID, action string) error
	DeleteServer(serverID string) error

	PatchServer(serverID string, definition ScalewayServerPatchDefinition) error
	ResolveServer(needle string) (ScalewayResolverResults, error)
	GetUserdatas(serverID string, metadata bool) (*ScalewayUserdatas, error)
	GetUserdata(serverID, key string, metadata bool) (*ScalewayUserdata, error)
	PatchUserdata(serverID, key string, value []byte, metadata bool) error
	DeleteUserdata(serverID, key string, metadata bool) error
}

func (s *ScalewayAPI) fetchServers(api string, query url.Values, out chan<- ScalewayServers) func() error {
	return func() error {
		resp, err := s.GetResponsePaginate(api, "servers", query)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
		if err != nil {
			return err
		}
		var servers ScalewayServers

		if err = json.Unmarshal(body, &servers); err != nil {
			return err
		}
		out <- servers
		return nil
	}
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
		panic("Not implemented yet")
	}
	if all && limit == 0 {
		s.Cache.ClearServers()
	}
	var (
		g    errgroup.Group
		apis = []string{
			ComputeAPIPar1,
			ComputeAPIAms1,
		}
	)

	serverChan := make(chan ScalewayServers, 2)
	for _, api := range apis {
		g.Go(s.fetchServers(api, query, serverChan))
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(serverChan)
	var servers ScalewayServers

	for server := range serverChan {
		servers.Servers = append(servers.Servers, server.Servers...)
	}

	for i, server := range servers.Servers {
		servers.Servers[i].DNSPublic = server.Identifier + URLPublicDNS
		servers.Servers[i].DNSPrivate = server.Identifier + URLPrivateDNS
		s.Cache.InsertServer(server.Identifier, server.Location.ZoneID, server.Arch, server.Organization, server.Name)
	}
	return &servers.Servers, nil
}

// GetServer gets a server from the ScalewayAPI
func (s *ScalewayAPI) GetServer(serverID string) (*ScalewayServer, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "servers/"+serverID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}

	var oneServer ScalewayOneServer

	if err = json.Unmarshal(body, &oneServer); err != nil {
		return nil, err
	}
	// FIXME arch, owner, title
	oneServer.Server.DNSPublic = oneServer.Server.Identifier + URLPublicDNS
	oneServer.Server.DNSPrivate = oneServer.Server.Identifier + URLPrivateDNS
	s.Cache.InsertServer(oneServer.Server.Identifier, oneServer.Server.Location.ZoneID, oneServer.Server.Arch, oneServer.Server.Organization, oneServer.Server.Name)
	return &oneServer.Server, nil
}

// PostServerAction posts an action on a server
func (s *ScalewayAPI) PostServerAction(serverID, action string) error {
	data := ScalewayServerAction{
		Action: action,
	}
	resp, err := s.PostResponse(s.computeAPI, fmt.Sprintf("servers/%s/action", serverID), data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusAccepted}, resp)
	return err
}

// DeleteServer deletes a server
func (s *ScalewayAPI) DeleteServer(serverID string) error {
	defer s.Cache.RemoveServer(serverID)
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("servers/%s", serverID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err = s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

// PostServer creates a new server
func (s *ScalewayAPI) PostServer(definition ScalewayServerDefinition) (string, error) {
	definition.Organization = s.Organization

	resp, err := s.PostResponse(s.computeAPI, "servers", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return "", err
	}
	var server ScalewayOneServer

	if err = json.Unmarshal(body, &server); err != nil {
		return "", err
	}
	// FIXME arch, owner, title
	s.Cache.InsertServer(server.Server.Identifier, server.Server.Location.ZoneID, server.Server.Arch, server.Server.Organization, server.Server.Name)
	return server.Server.Identifier, nil
}

// GetServerID returns exactly one server matching
func (s *ScalewayAPI) GetServerID(needle string) (string, error) {
	// Parses optional type prefix, i.e: "server:name" -> "name"
	_, needle = parseNeedle(needle)

	servers, err := s.ResolveServer(needle)
	if err != nil {
		return "", fmt.Errorf("Unable to resolve server %s: %s", needle, err)
	}
	if len(servers) == 1 {
		return servers[0].Identifier, nil
	}
	if len(servers) == 0 {
		return "", fmt.Errorf("No such server: %s", needle)
	}
	return "", showResolverResults(needle, servers)
}

// PatchServer updates a server
func (s *ScalewayAPI) PatchServer(serverID string, definition ScalewayServerPatchDefinition) error {
	resp, err := s.PatchResponse(s.computeAPI, fmt.Sprintf("servers/%s", serverID), definition)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusOK}, resp); err != nil {
		return err
	}
	return nil
}

// ResolveServer attempts to find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveServer(needle string) (ScalewayResolverResults, error) {
	servers, err := s.Cache.LookUpServers(needle, true)
	if err != nil {
		return servers, err
	}
	if len(servers) == 0 {
		if _, err = s.GetServers(true, 0); err != nil {
			return nil, err
		}
		servers, err = s.Cache.LookUpServers(needle, true)
	}
	return servers, err
}

// GetUserdatas gets list of userdata for a server
func (s *ScalewayAPI) GetUserdatas(serverID string, metadata bool) (*ScalewayUserdatas, error) {
	var uri, endpoint string

	endpoint = s.computeAPI
	if metadata {
		uri = "/user_data"
		endpoint = MetadataAPI
	} else {
		uri = fmt.Sprintf("servers/%s/user_data", serverID)
	}

	resp, err := s.GetResponsePaginate(endpoint, uri, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var userdatas ScalewayUserdatas

	if err = json.Unmarshal(body, &userdatas); err != nil {
		return nil, err
	}
	return &userdatas, nil
}

func (s *ScalewayUserdata) String() string {
	return string(*s)
}

// GetUserdata gets a specific userdata for a server
func (s *ScalewayAPI) GetUserdata(serverID, key string, metadata bool) (*ScalewayUserdata, error) {
	var uri, endpoint string

	endpoint = s.computeAPI
	if metadata {
		uri = fmt.Sprintf("/user_data/%s", key)
		endpoint = MetadataAPI
	} else {
		uri = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	var err error
	resp, err := s.GetResponsePaginate(endpoint, uri, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no such user_data %q (%d)", key, resp.StatusCode)
	}
	var data ScalewayUserdata
	data, err = ioutil.ReadAll(resp.Body)
	return &data, err
}

// PatchUserdata sets a user data
func (s *ScalewayAPI) PatchUserdata(serverID, key string, value []byte, metadata bool) error {
	var resource, endpoint string

	endpoint = s.computeAPI
	if metadata {
		resource = fmt.Sprintf("/user_data/%s", key)
		endpoint = MetadataAPI
	} else {
		resource = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	uri := fmt.Sprintf("%s/%s", strings.TrimRight(endpoint, "/"), resource)
	payload := new(bytes.Buffer)
	payload.Write(value)

	req, err := http.NewRequest("PATCH", uri, payload)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", s.Token)
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("User-Agent", s.userAgent)

	s.LogHTTP(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return fmt.Errorf("cannot set user_data (%d)", resp.StatusCode)
}

// DeleteUserdata deletes a server user_data
func (s *ScalewayAPI) DeleteUserdata(serverID, key string, metadata bool) error {
	var url, endpoint string

	endpoint = s.computeAPI
	if metadata {
		url = fmt.Sprintf("/user_data/%s", key)
		endpoint = MetadataAPI
	} else {
		url = fmt.Sprintf("servers/%s/user_data/%s", serverID, key)
	}

	resp, err := s.DeleteResponse(endpoint, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}
