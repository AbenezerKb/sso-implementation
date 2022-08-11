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
	// AccessToken is the access token for the current login
	AccessToken string `form:"access_token" query:"access_token" json:"access_token,omitempty"`
	// IDToken is the OpenID specific JWT token
	IDToken string `form:"id_token" query:"id_token" json:"id_token,omitempty"`
	// RefreshToken is the refresh token for the access token
	RefreshToken string `form:"refresh_token" query:"refresh_token" json:"refresh_token,omitempty"`
	// TokenType is the type of token
	TokenType string `form:"token_type" query:"token_type" json:"token_type,omitempty"`
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
