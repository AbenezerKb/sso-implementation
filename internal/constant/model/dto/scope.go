package dto

type Scope struct {
	// The scope name.
	Name string `json:"name,omitempty"`
	// The scope description.
	Description string `json:"description,omitempty"`
}
