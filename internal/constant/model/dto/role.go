package dto

// Role is a set of defined permissions that are grouped together with a name
type Role struct {
	// Name is a unique name for the role
	Name string `json:"name"`
	// Permissions are the list of permissions names this role contains
	Permissions []string `json:"permissions"`
}
