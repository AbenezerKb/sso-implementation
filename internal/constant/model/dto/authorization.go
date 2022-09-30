package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type Consent struct {
	AuthorizationRequestParam
	// ID is the unique id of this consent
	ID uuid.UUID `json:"id"`
	// Approved tells if this consent is approved by the user
	Approved bool `json:"approved"`
	// RequestOrigin is the origin of the client requesting authorization
	RequestOrigin string
}

type AuthCode struct {
	// The authorization code generated by the authorization server.
	Code string `json:"code"`
	// The client identifier.
	ClientID uuid.UUID `json:"client_id"`
	// The redirection URI used in the initial authorization request.
	RedirectURI string `json:"redirect_uri"`
	// The scope of the access request expressed as a list of space-delimited,
	Scope string `json:"scope"`
	// The state parameter passed in the initial authorization request.
	UserID uuid.UUID `json:"user_id"`
	// The state parameter passed in the initial authorization request.
	State string `json:"state"`
}

type AuthorizationRequestParam struct {
	// client identifier.
	ClientID uuid.UUID `form:"-" query:"-" json:"client_id,omitempty"`
	// redirection URI used in the initial authorization request.
	ResponseType string `form:"response_type" json:"response_type" query:"response_type"`
	// state parameter passed in the initial authorization request.
	State string `form:"state.omitempty" json:"state,omitempty" query:"state,omitempty"`
	// scope of the access request expressed as a list of space-delimited,
	Scope string `form:"scope" json:"scope" query:"scope"`
	// redirection URI used in the initial authorization request.
	RedirectURI string `form:"redirect_uri" json:"redirect_uri" query:"redirect_uri"`
	// specifies whether the Authorization Server MUST prompt the End-User for reauthentication.
	Prompt string `form:"prompt,omitempty" json:"prompt,omitempty" query:"prompt,omitempty"`
}

func (a *AuthorizationRequestParam) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ClientID, validation.Required.Error("client_id is required")),
		validation.Field(&a.ResponseType, validation.Required.Error("response_type is required"), validation.In("code", "token")),
		validation.Field(&a.Scope, validation.Required.Error("scope is required"), validation.In("openid", "profile", "email", "phone", "address", "offline_access")),
		validation.Field(&a.RedirectURI, validation.Required.Error("redirect_uri is required"), is.URL),
	)
}

type AuthHistory struct {
	// ID is the unique identifier for the history.
	// It is automatically generated when the history is created.
	ID uuid.UUID `json:"id"`
	// Code is the code the client uses to get access token.
	Code string `json:"code"`
	// UserID is the id of the user who granted access to the client.
	UserID uuid.UUID `json:"user_id"`
	// ClientID is the id of the client.
	ClientID uuid.UUID `json:"client_id"`
	// Scope is the scope the client is authorized to access.
	Scope string `json:"scope"`
	// Status is the status of the access token.
	// It can be either revoke or grant.
	Status string `json:"status"`
	// RedirectUri is the list of redirect uri of the client.
	RedirectUri string `json:"redirect_uri"`
	// CreatedAt is the time when the refresh token is created.
	// It is automatically set when the refresh token is created.
	CreatedAt time.Time `json:"created_at"`
}

type ConsentResponse struct {
	// Scopes is the list of scopes this consent holds
	Scopes []Scope `json:"scopes"`
	// ClientName is the name of the client
	ClientName string `json:"client_name"`
	// ClientLogo is the logo url of the client
	ClientLogo string `json:"client_logo"`
	// ClientType is the type of the client
	// It might be confidential or public
	ClientType string `json:"client_type"`
	// ClientTrusted tells if this client is a trusted first party client
	ClientTrusted bool `json:"client_trusted"`
	// ClientID is the id of the client given at the time of registration
	ClientID uuid.UUID `json:"client_id"`
	// UserID is the id of the user this consent is being given to
	UserID uuid.UUID `json:"user_id"`
	// Approved tells if this exact scope is previously approved by this user
	Approved bool `json:"approved"`
}

type LogoutRequest struct {
	IDTokenHint           string `form:"id_token_hint" json:"id_token_hint"`
	LogoutHint            string `form:"logout_hint" json:"logout_hint"`
	ClientID              string `form:"client_id" json:"client_id"`
	PostLogoutRedirectUri string `form:"post_logout_redirect_uri" json:"post_logout_redirect_uri"`
	State                 string `form:"state" json:"state"`
	UiLocales             string `form:"ui_locales" json:"ui_locales"`
}

func (l *LogoutRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.IDTokenHint, validation.Required.Error("login hint is required")),
		validation.Field(&l.PostLogoutRedirectUri, validation.Required.Error("post logout redirect uri is required")),
	)
}

type ConsentResultRsp struct {
	ConsentID     string `json:"consent_id"`
	FailureReason string `json:"failure_reason"`
}

// AuthorizedClientsResponse holds client data and access details for authorized client
type AuthorizedClientsResponse struct {
	Client
	// AuthGivenAt is the time this client is given authorization at
	AuthGivenAt time.Time `json:"created_at"`
	// AuthUpdatedAt is the time this authorization is last updated at
	AuthUpdatedAt time.Time `json:"updated_at"`
	// AuthExpiresAt is the time this authorization expires at
	AuthExpiresAt time.Time `json:"expires_at"`
	// AuthScopes is the scopes this authorization is given access to
	AuthScopes []Scope `json:"auth_scopes,omitempty"`
}
