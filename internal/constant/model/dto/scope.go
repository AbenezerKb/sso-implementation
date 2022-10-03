package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Scope struct {
	// The scope name.
	Name string `json:"name,omitempty"`
	// The scope description.
	Description string `json:"description,omitempty"`
	// resource server name
	ResourceServerName string `json:"resource_server_name,omitempty"`
	// date the scope created
	CreatedAt time.Time `json:"created_at"`
}

func (s Scope) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required.Error("name is required")),
		validation.Field(&s.Description, validation.Required.Error("description is required")),
	)
}
