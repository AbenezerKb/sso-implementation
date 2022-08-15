package dto

type Client struct {
	// ID is the unique identifier for the client.
	// It is automatically generated when the client is registered.
	ID string `json:"id"`
	// Name is the name of the client that will be displayed to the user.
	Name string `json:"name"`
	// ClientType is the type of the client.
	// It can be either confidential or public.
	ClientType string `json:"client_type"`
	// RedirectURIs is the list of redirect URIs of the client.
	RedirectURIs []string `json:"redirect_uris"`
	// Scopes is the list of default scopes of the client if one is not provided.
	Scopes string `json:"scopes"`
	// Secret is the secret the client uses to authenticate itself.
	// It is automatically generated when the client is registered.
	Secret string `json:"-"`
	// LogoURL is the URL of the client's logo.
	LogoURL string `json:"logo_url"`
}
