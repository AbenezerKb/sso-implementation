package dto

import (
	"fmt"
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
		validation.Field(&s.Name, validation.Required.Error("name is required"), validation.By(noSpaceScope)),
		validation.Field(&s.Description, validation.Required.Error("description is required")),
	)
}
func noSpaceScope(value interface{}) error {
	scope, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid value for scope")
	}
	scopeRunes := []rune(scope)
	for i := 0; i < len(scopeRunes); i++ {
		if scopeRunes[i] == ' ' {
			return fmt.Errorf("name can not contain spaces")
		}
	}
	return nil
}

type UpdateScopeParam struct {
	// The scope description.
	Description string `json:"description,omitempty"`
}

func (u UpdateScopeParam) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Description, validation.Required.Error("description is required")),
	)
}
