package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetSecurityGroups returns a ScalewaySecurityGroups
func (s *ScalewayAPI) GetSecurityGroups() (*ScalewayGetSecurityGroups, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "security_groups", url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups ScalewayGetSecurityGroups

	if err = json.Unmarshal(body, &securityGroups); err != nil {
		return nil, err
	}
	return &securityGroups, nil
}

// GetSecurityGroupRules returns a ScalewaySecurityGroupRules
func (s *ScalewayAPI) GetSecurityGroupRules(groupID string) (*ScalewayGetSecurityGroupRules, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("security_groups/%s/rules", groupID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroupRules ScalewayGetSecurityGroupRules

	if err = json.Unmarshal(body, &securityGroupRules); err != nil {
		return nil, err
	}
	return &securityGroupRules, nil
}

// GetASecurityGroupRule returns a ScalewaySecurityGroupRule
func (s *ScalewayAPI) GetASecurityGroupRule(groupID string, rulesID string) (*ScalewayGetSecurityGroupRule, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("security_groups/%s/rules/%s", groupID, rulesID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroupRules ScalewayGetSecurityGroupRule

	if err = json.Unmarshal(body, &securityGroupRules); err != nil {
		return nil, err
	}
	return &securityGroupRules, nil
}

// GetASecurityGroup returns a ScalewaySecurityGroup
func (s *ScalewayAPI) GetASecurityGroup(groupsID string) (*ScalewayGetSecurityGroup, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, fmt.Sprintf("security_groups/%s", groupsID), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var securityGroups ScalewayGetSecurityGroup

	if err = json.Unmarshal(body, &securityGroups); err != nil {
		return nil, err
	}
	return &securityGroups, nil
}

// PostSecurityGroup posts a group on a server
func (s *ScalewayAPI) PostSecurityGroup(group ScalewayNewSecurityGroup) error {
	resp, err := s.PostResponse(s.computeAPI, "security_groups", group)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusCreated}, resp)
	return err
}

// PostSecurityGroupRule posts a rule on a server
func (s *ScalewayAPI) PostSecurityGroupRule(SecurityGroupID string, rules ScalewayNewSecurityGroupRule) error {
	resp, err := s.PostResponse(s.computeAPI, fmt.Sprintf("security_groups/%s/rules", SecurityGroupID), rules)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusCreated}, resp)
	return err
}

// DeleteSecurityGroup deletes a SecurityGroup
func (s *ScalewayAPI) DeleteSecurityGroup(securityGroupID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("security_groups/%s", securityGroupID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}

// PutSecurityGroup updates a SecurityGroup
func (s *ScalewayAPI) PutSecurityGroup(group ScalewayUpdateSecurityGroup, securityGroupID string) error {
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("security_groups/%s", securityGroupID), group)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// PutSecurityGroupRule updates a SecurityGroupRule
func (s *ScalewayAPI) PutSecurityGroupRule(rules ScalewayNewSecurityGroupRule, securityGroupID, RuleID string) error {
	resp, err := s.PutResponse(s.computeAPI, fmt.Sprintf("security_groups/%s/rules/%s", securityGroupID, RuleID), rules)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// DeleteSecurityGroupRule deletes a SecurityGroupRule
func (s *ScalewayAPI) DeleteSecurityGroupRule(SecurityGroupID, RuleID string) error {
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("security_groups/%s/rules/%s", SecurityGroupID, RuleID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}
