package dto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"time"
)

// ResourceServer is a server that this sso controls access for
type ResourceServer struct {
	// ID is the unique id for this resource server
	ID uuid.UUID `json:"id"`
	// Name is the resource server name.
	// It must be unique across the sso
	Name string `json:"name"`
	// CreatedAt is the time this resource server is created at
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time this resource server is updated at
	UpdatedAt time.Time `json:"updated_at"`
	// Scopes is the scopes of this resource server
	Scopes []Scope `json:"scopes,omitempty"`
	// Secret is the secret of the resource server that will be used for authentication on the sso.
	Secret string `json:"secret,omitempty"`
}

func (r ResourceServer) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("server name is required")),
		validation.Field(&r.Scopes, validation.By(scopesValidate)),
	)
}

func scopesValidate(value interface{}) error {
	scopes, ok := value.([]Scope)
	if !ok {
		return fmt.Errorf("invalid scopes")
	}

	for i := 0; i < len(scopes); i++ {
		if err := validation.Validate(scopes[i].Name, validation.Required.Error("scope name is required")); err != nil {
			return err
		}
		if err := validation.Validate(scopes[i].Description, validation.Required.Error("scope description is required")); err != nil {
			return err
		}
		for j := 0; j < len(scopes); j++ {
			if scopes[i].Name == scopes[j].Name && i != j {
				return fmt.Errorf("scope name must be unique")
			}
		}
	}

	return nil
}
