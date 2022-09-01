package request_models

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RevokeClientBody struct {
	// ClientID is the id of the client to be revoked access of.
	// It is a required field
	ClientID string `json:"client_id"`
}

func (r RevokeClientBody) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ClientID, validation.Required.Error("client_id is required"), is.UUID.Error("invalid client_id")))
}
