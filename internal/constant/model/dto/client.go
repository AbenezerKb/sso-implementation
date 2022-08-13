package dto

type Client struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	ClientType   string   `json:"client_type"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       string   `json:"scopes"`
	Secret       string   `json:"-"`
	LogoURL      string   `json:"logo_url"`
}
