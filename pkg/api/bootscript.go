package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ResolveBootscript attempts to find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveBootscript(needle string) (ScalewayResolverResults, error) {
	bootscripts, err := s.Cache.LookUpBootscripts(needle, true)
	if err != nil {
		return bootscripts, err
	}
	if len(bootscripts) == 0 {
		if _, err = s.GetBootscripts(); err != nil {
			return nil, err
		}
		bootscripts, err = s.Cache.LookUpBootscripts(needle, true)
	}
	return bootscripts, err
}

// GetBootscripts gets the list of bootscripts from the ScalewayAPI
func (s *ScalewayAPI) GetBootscripts() (*[]ScalewayBootscript, error) {
	query := url.Values{}

	s.Cache.ClearBootscripts()
	resp, err := s.GetResponsePaginate(s.computeAPI, "bootscripts", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var bootscripts ScalewayBootscripts

	if err = json.Unmarshal(body, &bootscripts); err != nil {
		return nil, err
	}
	for _, bootscript := range bootscripts.Bootscripts {
		// FIXME region, arch, owner, title
		s.Cache.InsertBootscript(bootscript.Identifier, "", bootscript.Arch, bootscript.Organization, bootscript.Title)
	}
	return &bootscripts.Bootscripts, nil
}

// GetBootscript gets a bootscript from the ScalewayAPI
func (s *ScalewayAPI) GetBootscript(bootscriptID string) (*ScalewayBootscript, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "bootscripts/"+bootscriptID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var oneBootscript ScalewayOneBootscript

	if err = json.Unmarshal(body, &oneBootscript); err != nil {
		return nil, err
	}
	// FIXME region, arch, owner, title
	s.Cache.InsertBootscript(oneBootscript.Bootscript.Identifier, "", oneBootscript.Bootscript.Arch, oneBootscript.Bootscript.Organization, oneBootscript.Bootscript.Title)
	return &oneBootscript.Bootscript, nil
}

// GetBootscriptID returns exactly one bootscript matching
func (s *ScalewayAPI) GetBootscriptID(needle, arch string) (string, error) {
	// Parses optional type prefix, i.e: "bootscript:name" -> "name"
	_, needle = parseNeedle(needle)

	bootscripts, err := s.ResolveBootscript(needle)
	if err != nil {
		return "", fmt.Errorf("Unable to resolve bootscript %s: %s", needle, err)
	}
	bootscripts.FilterByArch(arch)
	if len(bootscripts) == 1 {
		return bootscripts[0].Identifier, nil
	}
	if len(bootscripts) == 0 {
		return "", fmt.Errorf("No such bootscript: %s", needle)
	}
	return "", showResolverResults(needle, bootscripts)
}
