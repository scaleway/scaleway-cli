package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type IPAPI interface {
	GetIPS() (*ScalewayGetIPS, error)
	NewIP() (*ScalewayGetIP, error)
	AttachIP(ipID, serverID string) error
	DetachIP(ipID string) error
	DeleteIP(ipID string) error
	GetIP(ipID string) (*ScalewayGetIP, error)
}

// GetIPS returns a ScalewayGetIPS
func (s *ScalewayAPI) GetIPS() (*ScalewayGetIPS, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "ips", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var ips ScalewayGetIPS

	if err = json.Unmarshal(body, &ips); err != nil {
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
	resp, err := s.PostResponse(s.computeAPI, "ips", orga)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return nil, err
	}
	var ip ScalewayGetIP

	if err = json.Unmarshal(body, &ip); err != nil {
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
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("ips/%s", ipID), update)
	if err != nil {
		return err
	}
	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// DetachIP detaches an IP from a server
func (s *ScalewayAPI) DetachIP(ipID string) error {
	ip, err := s.GetIP(ipID)
	if err != nil {
		return err
	}
	ip.IP.Server = nil
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("ips/%s", ipID), ip.IP)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// DeleteIP deletes an IP
func (s *ScalewayAPI) DeleteIP(ipID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("ips/%s", ipID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}

// GetIP returns a ScalewayGetIP
func (s *ScalewayAPI) GetIP(ipID string) (*ScalewayGetIP, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("ips/%s", ipID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var ip ScalewayGetIP

	if err = json.Unmarshal(body, &ip); err != nil {
		return nil, err
	}
	return &ip, nil
}
