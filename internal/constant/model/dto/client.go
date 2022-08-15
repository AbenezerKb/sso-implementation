package dto

import (
	"time"
)

type Client struct {
	ID     string `json:"id,omitempty"`
	Secret string `json:"secret,omitempty"`
	//array
	RedirectURIs []string `json:"redirect_uris,omitempty"`

	GrantTypes string `json:"grant_types,omitempty"`
	Name       string `json:"name,omitempty"`

	LogoURI string `json:"logo_uri,omitempty"`
	Scope   string `json:"scope,omitempty"`

	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"-"`
}
