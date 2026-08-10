//go:build darwin || linux || windows

package object

import "encoding/json"

// StringList is a simple slice of string with custom unmarshalling methods
type StringList []string

func (s *StringList) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = StringList{single}
		return nil
	}

	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		return err
	}

	*s = multi
	return nil
}

type BucketPolicy struct {
	Version   string                   `json:"Version"`
	ID        string                   `json:"Id,omitempty"`
	Statement []*BucketPolicyStatement `json:"Statement,omitempty"`
}

type BucketPolicyStatement struct {
	Sid       string                    `json:"Sid,omitempty"`
	Effect    string                    `json:"Effect"`
	Principal BucketPolicyPrincipal     `json:"Principal"`
	Action    StringList                `json:"Action"`
	Resource  StringList                `json:"Resource"`
	Condition map[string]map[string]any `json:"Condition,omitempty"`
}

type BucketPolicyPrincipal struct {
	SCW StringList `json:"SCW,omitempty"`
	Raw string     `json:",omitempty"` // This interpolates the "*" notation
}

func (p *BucketPolicyPrincipal) UnmarshalJSON(data []byte) error {
	var wildcard string
	if err := json.Unmarshal(data, &wildcard); err == nil {
		p.Raw = wildcard
		return nil
	}

	type Alias BucketPolicyPrincipal
	var temp Alias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*p = BucketPolicyPrincipal(temp)

	return nil
}

func (p *BucketPolicyPrincipal) MarshalJSON() ([]byte, error) {
	if p.Raw != "" {
		return json.Marshal(p.Raw)
	}

	type Alias BucketPolicyPrincipal
	return json.Marshal((Alias)(*p))
}
