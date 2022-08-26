package dto

import (
	"sso/internal/constant"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
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
	// ExpiresAt is time the access token is going to be expired.
	ExpiresIn string `json:"expires_in"`
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

type AccessTokenRequest struct {
	// GrantType is the type of flow the client is following to get access token.
	// It can be either authorization_code or refresh_token.
	GrantType string `json:"grant_type" form:"grant_type"`
	// Authorization code generated by the authorization server.
	Code string `json:"code" form:"code"`
	// Redirection URI used in the initial authorization request.
	RedirectURI string `json:"redirect_uri" form:"redirect_uri"`
	// RefreshToken is the opaque string that was given by the auth server when issuing the access token.
	// it's used to refresh the access token.
	RefreshToken string `json:"refresh_token"`
}

func (a AccessTokenRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Code, validation.When(a.GrantType == constant.AuthorizationCode, validation.Required.Error("code is required"))),
		validation.Field(&a.RedirectURI, validation.When(a.GrantType == constant.AuthorizationCode, validation.Required.Error("redirect_uri is required")), is.URL.Error("invalid redirect_uri")),
		validation.Field(&a.GrantType, validation.Required.Error("grant_type is required"), validation.In(constant.AuthorizationCode, constant.RefreshToken)),
		validation.Field(&a.RefreshToken, validation.When(a.GrantType == constant.RefreshToken, validation.Required.Error("refresh_token is required"))),
	)
}

type RefreshToken struct {
	// ID is the unique identifier for the refresh token.
	// It is automatically generated when the refresh token is created.
	ID uuid.UUID `json:"id"`
	// Refreshtoken is the opaque string the client uses to refresh access token.
	Refreshtoken string `json:"refreshtoken"`
	// Code is the code the client uses to get access token.
	Code string `json:"code"`
	// UserID is the id of the user who granted access to the client.
	UserID uuid.UUID `json:"user_id"`
	// ClientID is the id of the client.
	ClientID uuid.UUID `json:"client_id"`
	// Scope is the scope the client is authorized to access.
	Scope string `json:"scope"`
	// RedirectUri is the list of redirect uri of the client.
	RedirectUri string `json:"redirect_uri"`
	// ExpiresAt is time the refresh token is going to be expired.
	ExpiresAt time.Time `json:"expires_at"`
	// CreatedAt is the time when the refresh token is created.
	// It is automatically set when the refresh token is created.
	CreatedAt time.Time `json:"created_at"`
}

type InternalRefreshToken struct {
	// ID is the unique identifier for the refresh token.
	// It is automatically generated when the refresh token is created.
	ID uuid.UUID `json:"id"`
	// Refreshtoken is the opaque string users uses to refresh access token.
	Refreshtoken string `json:"refreshtoken"`
	// ExpiresAt is time the refresh token is going to be expired.
	// UserID is the id of the user who granted access to the client.
	UserID uuid.UUID `json:"user_id"`
	// ExpiresAt is time the refresh token is going to be expired.
	ExpiresAt time.Time `json:"expires_at"`
	// CreatedAt is the time when the refresh token is created.
	// It is automatically set when the refresh token is created.
	CreatedAt time.Time `json:"created_at"`
}
