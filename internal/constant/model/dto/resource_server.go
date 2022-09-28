package dto

// ResourceServer is a server that this sso controls access for
type ResourceServer struct {
	// Name is the resource server name.
	// It must be unique across the sso
	Name string `json:"name"`
}
