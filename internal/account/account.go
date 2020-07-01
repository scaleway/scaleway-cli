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
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	OrganizationID string `json:"organization_id"`
	ProjectID      string `json:"project_id"`
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

	resp, err := extractHTTPClient(ctx).Post(accountURL+"/tokens", "application/json", bytes.NewReader(rawJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return &LoginResponse{
			TwoFactorRequired: true,
		}, nil
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return &LoginResponse{
			WrongPassword: true,
		}, nil
	}
	if resp.StatusCode != http.StatusCreated {
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

func GetAPIKey(ctx context.Context, secretKey string) (*Token, error) {
	resp, err := extractHTTPClient(ctx).Get(accountURL + "/tokens/" + secretKey)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not get token")
	}

	token := &LoginResponse{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, token)
	if err != nil {
		return nil, err
	}

	return token.Token, err
}
