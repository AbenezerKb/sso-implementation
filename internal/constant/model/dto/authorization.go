package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Consent struct {
	AuthorizationRequestParam
	ID string `json:"id"`
	// The consent status.
	Approved bool `json:"approved"`
	// Users Id
	UserID string `json:"userID"`
	// Roles of the user.
	Roles string `json:"roles"`
}

type AuthCode struct {
	// The authorization code generated by the authorization server.
	Code string `json:"code"`
	// The client identifier.
	ClientID string `json:"client_id"`
	// The redirection URI used in the initial authorization request.
	RedirectURI string `json:"redirect_uri"`
	// The scope of the access request expressed as a list of space-delimited,
	Scope string `json:"scope"`
	// The state parameter passed in the initial authorization request.
	UserID string `json:"user_id"`
	// The state parameter passed in the initial authorization request.
	State string `json:"state"`
}

type ConsentData struct {
	Consent
	// The user data
	User *User `json:"user"`
	// The client data
	Client *Client `json:"client"`
	// The scope data
	Scopes []Scope `json:"scopes"`
}

type AuthorizationRequestParam struct {
	// The client identifier.
	ClientID string `form:"client_id" query:"client_id" json:"client_id,omitempty"`
	// The redirection URI used in the initial authorization request.
	ResponseType string `form:"response_type" json:"response_type" query:"response_type"`
	// the state parameter passed in the initial authorization request.
	State string `form:"state" json:"state" query:"state"`
	// The scope of the access request expressed as a list of space-delimited,
	Scope string `form:"scope" json:"scope" query:"scope"`
	// The redirection URI used in the initial authorization request.
	RedirectURI string `form:"redirect_uri" json:"redirect_uri" query:"redirect_uri"`
}

func (a *AuthorizationRequestParam) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.ClientID, validation.Required.Error("client_id is required")),
		validation.Field(&a.ResponseType, validation.Required.Error("response_type is required"), validation.In("code", "token")),
		validation.Field(&a.State, validation.Required.Error("state is required")),
		validation.Field(&a.Scope, validation.Required.Error("scope is required"), validation.In("openid", "profile", "email", "phone", "address", "offline_access")),
		validation.Field(&a.RedirectURI, validation.Required.Error("redirect_uri is required"), is.URL),
	)
}
