package request_models

import validation "github.com/go-ozzo/ozzo-validation/v4"

// LoginWithIP is used to request login with external identity providers
type LoginWithIP struct {
	// Code is the authorization code from the identity provider
	Code string `json:"code"`
	// IdentityProvider is the id of the identity provider to log in with
	IdentityProvider string `json:"ip"`
}

func (l LoginWithIP) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Code, validation.Required.Error("code is required")),
		validation.Field(&l.IdentityProvider, validation.Required.Error("identity provider is required")),
	)
}
