package account

import (
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
