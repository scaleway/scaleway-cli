package account

import (
	"bytes"
	"context"
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
	Token             *Token `json:"token"`
	TwoFactorRequired bool   `json:"-"`
	WrongPassword     bool   `json:"-"`
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
)

// Login creates a new token
func Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// todo: add line of log
	rawJSON, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := extractHttpClient(ctx).Post(accountURL+"/tokens", "application/json", bytes.NewReader(rawJSON))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 403 {
		return &LoginResponse{
			TwoFactorRequired: true,
		}, nil
	}
	if resp.StatusCode == 401 {
		return &LoginResponse{
			WrongPassword: true,
		}, nil
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("scaleway-cli: %s", resp.Status)
	}
	loginResponse := &LoginResponse{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, loginResponse)
	return loginResponse, err
}

func GetAccessKey(ctx context.Context, secretKey string) (string, error) {
	resp, err := extractHttpClient(ctx).Get(accountURL + "/tokens/" + secretKey)
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

func getOrganizations(ctx context.Context, secretKey string) ([]organization, error) {
	req, err := http.NewRequest("GET", accountURL+"/organizations", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Auth-Token", secretKey)
	resp, err := extractHttpClient(ctx).Do(req)
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

func GetOrganizationsIds(ctx context.Context, secretKey string) ([]string, error) {
	organizations, err := getOrganizations(ctx, secretKey)
	if err != nil {
		return nil, err
	}
	ids := []string(nil)
	for _, organization := range organizations {
		ids = append(ids, organization.ID)
	}
	return ids, nil
}
