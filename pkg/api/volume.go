package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DeleteVolume deletes a volume
func (s *ScalewayAPI) DeleteVolume(volumeID string) error {
	defer s.Cache.RemoveVolume(volumeID)
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("volumes/%s", volumeID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

// GetVolumeID returns exactly one volume matching
func (s *ScalewayAPI) GetVolumeID(needle string) (string, error) {
	// Parses optional type prefix, i.e: "volume:name" -> "name"
	_, needle = parseNeedle(needle)

	volumes, err := s.ResolveVolume(needle)
	if err != nil {
		return "", fmt.Errorf("Unable to resolve volume %s: %s", needle, err)
	}
	if len(volumes) == 1 {
		return volumes[0].Identifier, nil
	}
	if len(volumes) == 0 {
		return "", fmt.Errorf("No such volume: %s", needle)
	}
	return "", showResolverResults(needle, volumes)
}

// GetVolumes gets the list of volumes from the ScalewayAPI
func (s *ScalewayAPI) GetVolumes() (*[]ScalewayVolume, error) {
	query := url.Values{}
	s.Cache.ClearVolumes()

	resp, err := s.GetResponsePaginate(s.computeAPI, "volumes", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}

	var volumes ScalewayVolumes

	if err = json.Unmarshal(body, &volumes); err != nil {
		return nil, err
	}
	for _, volume := range volumes.Volumes {
		// FIXME region, arch, owner, title
		s.Cache.InsertVolume(volume.Identifier, "", "", volume.Organization, volume.Name)
	}
	return &volumes.Volumes, nil
}

// GetVolume gets a volume from the ScalewayAPI
func (s *ScalewayAPI) GetVolume(volumeID string) (*ScalewayVolume, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "volumes/"+volumeID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var oneVolume ScalewayOneVolume

	if err = json.Unmarshal(body, &oneVolume); err != nil {
		return nil, err
	}
	// FIXME region, arch, owner, title
	s.Cache.InsertVolume(oneVolume.Volume.Identifier, "", "", oneVolume.Volume.Organization, oneVolume.Volume.Name)
	return &oneVolume.Volume, nil
}

// PostVolume creates a new volume
func (s *ScalewayAPI) PostVolume(definition ScalewayVolumeDefinition) (string, error) {
	definition.Organization = s.Organization
	if definition.Type == "" {
		definition.Type = "l_ssd"
	}

	resp, err := s.PostResponse(s.computeAPI, "volumes", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return "", err
	}
	var volume ScalewayOneVolume

	if err = json.Unmarshal(body, &volume); err != nil {
		return "", err
	}
	// FIXME: s.Cache.InsertVolume(volume.Volume.Identifier, volume.Volume.Name)
	return volume.Volume.Identifier, nil
}

// PutVolume updates a volume
func (s *ScalewayAPI) PutVolume(volumeID string, definition ScalewayVolumePutDefinition) error {
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("volumes/%s", volumeID), definition)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// ResolveVolume attempts to find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveVolume(needle string) (ScalewayResolverResults, error) {
	volumes, err := s.Cache.LookUpVolumes(needle, true)
	if err != nil {
		return volumes, err
	}
	if len(volumes) == 0 {
		if _, err = s.GetVolumes(); err != nil {
			return nil, err
		}
		volumes, err = s.Cache.LookUpVolumes(needle, true)
	}
	return volumes, err
}
