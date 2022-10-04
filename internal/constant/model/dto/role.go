package dto

import (
	"sso/internal/constant/permissions"
)

// Role is a set of defined permissions that are grouped together with a name
type Role struct {
	// Name is a unique name for the role
	Name string `json:"name"`
	// Permissions are the list of permissions this role contains
	Permissions []permissions.Permission `json:"permissions"`
}
