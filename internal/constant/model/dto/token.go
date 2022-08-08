package dto

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessToken struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
	Phone      string `json:"phone"`

	ClientID  string     `form:"client_id" query:"client_id" json:"client_id,omitempty"`
	UserID    string     `form:"user_id" query:"user_id" json:"user_id,omitempty"`
	Roles     string     `form:"roles" query:"roles" json:"roles,omitempty"`
	Scope     string     `form:"scope" query:"scope" json:"scope,omitempty"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	AccessToken  string `form:"access_token" query:"access_token" json:"access_token,omitempty"`
	IDToken      string `form:"id_token" query:"id_token" json:"id_token,omitempty"`
	RefreshToken string `form:"refresh_token" query:"refresh_token" json:"refresh_token,omitempty"`
	TokenType    string `form:"token_type" query:"token_type" json:"token_type,omitempty"`

	Issuer         string `json:"iss,omitempty"`
	Subject        string `json:"sub,omitempty"`
	Audience       string `json:"aud,omitempty"`
	Expiry         int64  `json:"expiry,omitempty"`
	NotBefore      int64  `json:"nbf,omitempty"`
	IssuedAt       int64  `json:"iat,omitempty"`
	ID             string `json:"jti,omitempty"`
	PasswordStatus string `json:"password_status,omitempty"`
	ExpiresIn      int64  `form:"expires_in" query:"expires_in" json:"expires_in,omitempty"`
	Scope          string `form:"scope" query:"scope" json:"scope,omitempty"`
	State          string `form:"state" query:"scstateope" json:"state,omitempty"`
	Roles          string `json:"roles,omitempty"`
}

type IDTokenPayload struct {
	FirstName       string `json:"first_name"`
	MiddleName      string `json:"middle_name"`
	LastName        string `json:"last_name"`
	Picture         string `json:"picture"`
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone"`
	AuthorizedParty string `json:"azp"`

	jwt.RegisteredClaims
}
