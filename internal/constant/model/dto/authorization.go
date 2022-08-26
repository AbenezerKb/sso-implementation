package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type Consent struct {
	AuthorizationRequestParam
	ID uuid.UUID `json:"id"`
	// The consent status.
	Approved bool `json:"approved"`
	// Users Id
	UserID uuid.UUID `json:"userID"`
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
	Scopes []Scope `json:"scopes"`
	// client name
	Client_Name string `json:"client_name"`
	// client logo
	Client_Logo string `json:"client_logo"`
	// client Type
	Client_Type string `json:"client_type"`
	// is the client fully trusted
	Client_Trusted bool `json:"client_trusted"`
	// client id
	Client_ID uuid.UUID `json:"client_id"`
	// user id
	User_ID uuid.UUID `json:"user_id"`
	// whether the user has consented to the scopes
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
