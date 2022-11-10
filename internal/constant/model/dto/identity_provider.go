package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"time"
)

// IdentityProvider is an authorization server that supports openid connect.
// This holds the client information that is registered on this identity provider's server.
type IdentityProvider struct {
	// ID is the id of this identity provider
	ID uuid.UUID `json:"id"`
	// Name is the name of this identity provider.
	Name string `json:"name"`
	// LogoURI is the uri of a logo that will be shown for this identity provider
	LogoURI string `json:"logo_uri,omitempty"`
	// ClientID is the id of the client that is registered on the server of this identity provider.
	// Requests to this identity provider are passed on behalf of this client id.
	ClientID string `json:"client_id"`
	// ClientSecret is the password to be used with the ClientID
	ClientSecret string `json:"client_secret"`
	// RedirectURI is the redirect uri the client with ClientID has registered on the server of this identity provider.
	RedirectURI string `json:"redirect_uri"`
	// AuthorizationURI is the uri to request openid authorization from
	AuthorizationURI string `json:"authorization_uri"`
	// TokenEndpointURI is the uri to exchange code with access token
	TokenEndpointURI string `json:"token_endpoint_uri"`
	// UserInfoEndpointURI is the uri to exchange access token with user profile information
	UserInfoEndpointURI string `json:"user_info_endpoint_uri,omitempty"`
	// Status is the status of this identity provider
	Status string `json:"status"`
	// CreatedAt is the time this identity provider was created at
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time this identity provider was last updated at
	UpdatedAt time.Time `json:"updated_at"`
}

func (i IdentityProvider) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Name, validation.Required.Error("name is required")),
		validation.Field(&i.LogoURI, is.URL.Error("invalid logo_uri"), validation.By(ValidateLogo)),
		validation.Field(&i.ClientID, validation.Required.Error("client_id is required")),
		validation.Field(&i.ClientSecret, validation.Required.Error("client_secret is required")),
		validation.Field(&i.RedirectURI, validation.Required.Error("redirect_uri is required")),
		validation.Field(&i.AuthorizationURI, validation.Required.Error("authorization_uri is required"), is.URL.Error("invalid authorization_uri")),
		validation.Field(&i.TokenEndpointURI, validation.Required.Error("token_endpoint_uri is required"), is.URL.Error("invalid token_endpoint_uri")),
		validation.Field(&i.UserInfoEndpointURI, is.URL.Error("invalid user_info_endpoint_uri")),
	)
}
