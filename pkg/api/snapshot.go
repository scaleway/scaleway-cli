package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SnapshotAPI interface {
	DeleteSnapshot(snapshotID string) error
	GetSnapshots() (*[]ScalewaySnapshot, error)
	GetSnapshot(snapshotID string) (*ScalewaySnapshot, error)
	PostSnapshot(volumeID string, name string) (string, error)
}

// DeleteSnapshot deletes a snapshot
func (s *ScalewayAPI) DeleteSnapshot(snapshotID string) error {
	defer s.Cache.RemoveSnapshot(snapshotID)
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("snapshots/%s", snapshotID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

// GetSnapshots gets the list of snapshots from the ScalewayAPI
func (s *ScalewayAPI) GetSnapshots() (*[]ScalewaySnapshot, error) {
	query := url.Values{}
	s.Cache.ClearSnapshots()

	resp, err := s.GetResponsePaginate(s.computeAPI, "snapshots", query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var snapshots ScalewaySnapshots

	if err = json.Unmarshal(body, &snapshots); err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots.Snapshots {
		// FIXME region, arch, owner, title
		s.Cache.InsertSnapshot(snapshot.Identifier, "", "", snapshot.Organization, snapshot.Name)
	}
	return &snapshots.Snapshots, nil
}

// GetSnapshot gets a snapshot from the ScalewayAPI
func (s *ScalewayAPI) GetSnapshot(snapshotID string) (*ScalewaySnapshot, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "snapshots/"+snapshotID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var oneSnapshot ScalewayOneSnapshot

	if err = json.Unmarshal(body, &oneSnapshot); err != nil {
		return nil, err
	}
	// FIXME region, arch, owner, title
	s.Cache.InsertSnapshot(oneSnapshot.Snapshot.Identifier, "", "", oneSnapshot.Snapshot.Organization, oneSnapshot.Snapshot.Name)
	return &oneSnapshot.Snapshot, nil
}

// PostSnapshot creates a new snapshot
func (s *ScalewayAPI) PostSnapshot(volumeID string, name string) (string, error) {
	definition := ScalewaySnapshotDefinition{
		VolumeIDentifier: volumeID,
		Name:             name,
		Organization:     s.Organization,
	}
	resp, err := s.PostResponse(s.computeAPI, "snapshots", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return "", err
	}
	var snapshot ScalewayOneSnapshot

	if err = json.Unmarshal(body, &snapshot); err != nil {
		return "", err
	}
	// FIXME arch, owner, title
	s.Cache.InsertSnapshot(snapshot.Snapshot.Identifier, "", "", snapshot.Snapshot.Organization, snapshot.Snapshot.Name)
	return snapshot.Snapshot.Identifier, nil
}

// GetSnapshotID returns exactly one snapshot matching
func (s *ScalewayAPI) GetSnapshotID(needle string) (string, error) {
	// Parses optional type prefix, i.e: "snapshot:name" -> "name"
	_, needle = parseNeedle(needle)

	snapshots, err := s.ResolveSnapshot(needle)
	if err != nil {
		return "", fmt.Errorf("Unable to resolve snapshot %s: %s", needle, err)
	}
	if len(snapshots) == 1 {
		return snapshots[0].Identifier, nil
	}
	if len(snapshots) == 0 {
		return "", fmt.Errorf("No such snapshot: %s", needle)
	}
	return "", showResolverResults(needle, snapshots)
}

// ResolveSnapshot attempts to find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveSnapshot(needle string) (ScalewayResolverResults, error) {
	snapshots, err := s.Cache.LookUpSnapshots(needle, true)
	if err != nil {
		return snapshots, err
	}
	if len(snapshots) == 0 {
		if _, err = s.GetSnapshots(); err != nil {
			return nil, err
		}
		snapshots, err = s.Cache.LookUpSnapshots(needle, true)
	}
	return snapshots, err
}
