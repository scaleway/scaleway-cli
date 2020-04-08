package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Token represents a Token
type Token struct {
	UserID    string `json:"user_id"`
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	ID        string `json:"id"`
}

type LoginResponse struct {
	Token *Token `json:"token"`
}

type LoginRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password,omitempty"`
	TwoFactorToken string `json:"2FA_token,omitempty"`
	Description    string `json:"description,omitempty"`
	Expires        bool   `json:"expires"`
}

type organizationsResponse struct {
	Organizations []organization `json:"organizations"`
}

type organization struct {
	ID string `json:"id"`
}

var (
	accountURL = "https://account.scaleway.com"

	WrongPassword = fmt.Errorf("wrong password")
)

// Login creates a new token
func Login(req *LoginRequest) (t *Token, twoFactorRequired bool, err error) {
	// todo: add line of log
	rawJSON, err := json.Marshal(req)
	if err != nil {
		return nil, false, err
	}

	resp, err := http.Post(accountURL+"/tokens", "application/json", bytes.NewReader(rawJSON))
	if err != nil {
		return nil, false, err
	}

	if resp.StatusCode == 403 {
		return nil, true, nil
	}
	if resp.StatusCode == 401 {
		return nil, false, WrongPassword
	}
	if resp.StatusCode != 201 {
		return nil, false, fmt.Errorf("scaleway-cli: %s", resp.Status)
	}
	token := &LoginResponse{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	err = json.Unmarshal(b, token)
	return token.Token, false, err
}

func GetAccessKey(secretKey string) (string, error) {
	resp, err := http.Get(accountURL + "/tokens/" + secretKey)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not get token")
	}

	token := &LoginResponse{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(b, token)
	if err != nil {
		return "", err
	}

	return token.Token.AccessKey, err
}

func getOrganizations(secretKey string) ([]organization, error) {
	req, err := http.NewRequest("GET", accountURL+"/organizations", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", secretKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get organizations from %s", accountURL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	organizationsResponse := &organizationsResponse{}
	err = json.Unmarshal(body, organizationsResponse)
	if err != nil {
		return nil, err
	}
	return organizationsResponse.Organizations, nil
}

func GetOrganizationsIds(secretKey string) ([]string, error) {
	organizations, err := getOrganizations(secretKey)
	if err != nil {
		return nil, err
	}
	ids := []string(nil)
	for _, organization := range organizations {
		ids = append(ids, organization.ID)
	}
	return ids, nil
}
