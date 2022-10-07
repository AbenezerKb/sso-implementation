package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"sso/internal/constant"
	"time"
)

// Role is a set of defined permissions that are grouped together with a name
type Role struct {
	// Name is a unique name for the role
	Name string `json:"name"`
	// Permissions are the list of permissions names this role contains
	Permissions []string `json:"permissions"`
	// Status is the current status of this role
	Status string `json:"status"`
	// CreatedAt is the time this role is created on
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time this role is last updated at
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Role) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required.Error("name is required")),
		validation.Field(&r.Permissions, validation.Required.Error("permissions is required")))
}

type UpdateRoleStatus struct {
	// Status is new status that will replace old status of the role
	Status string `json:"status"`
}

func (u UpdateRoleStatus) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Status, validation.Required.Error("status is required"), validation.In(constant.Active, constant.Inactive, constant.Pending).Error("invalid status")),
	)
}

type UpdateRole struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func (u UpdateRole) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required.Error("name is required")),
		validation.Field(&u.Permissions, validation.Required.Error("permissions is required")),
	)
}
