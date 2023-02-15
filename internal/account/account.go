package account

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
