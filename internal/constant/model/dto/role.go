package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

// Role is a set of defined permissions that are grouped together with a name
type Role struct {
	// Name is a unique name for the role
	Name string `json:"name"`
	// Permissions are the list of permissions names this role contains
	Permissions []string `json:"permissions"`
}

func (r Role) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("name is required")),
		validation.Field(&r.Permissions, validation.Required.Error("permissions is required")))
}
