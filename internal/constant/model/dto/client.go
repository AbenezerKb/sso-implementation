package dto

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"net/http"
	"strings"
)

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
	// Each redirect URI must be a valid URL and must use HTTPS.
	RedirectURIs []string `json:"redirect_uris"`
	// Scopes is the list of default scopes of the client if one is not provided.
	Scopes string `json:"scopes"`
	// Secret is the secret the client uses to authenticate itself.
	// It is automatically generated when the client is registered.
	Secret string `json:"-"`
	// LogoURL is the URL of the client's logo.
	// It must be a valid URL.
	LogoURL string `json:"logo_url"`
}

func (c Client) ValidateClient() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required.Error("name is required"), validation.Length(3, 32).Error("name must be between 3 and 32 characters")),
		validation.Field(&c.ClientType, validation.Required.Error("client type is required"), validation.In("confidential", "public").Error("client type must be either confidential or public")),
		validation.Field(&c.RedirectURIs, validation.Required.Error("redirect URIs are required"), validation.Each(is.URL.Error("redirect URI is not valid"), validation.By(ValidateURI))),
		validation.Field(&c.Scopes, validation.Required.Error("scopes is required")),
		validation.Field(&c.LogoURL, validation.Required.Error("logo URL is required"), is.URL.Error("logo URL is not valid"), validation.By(ValidateURI)),
	)

}

func ValidateURI(uri interface{}) error {
	if !strings.HasPrefix(fmt.Sprint(uri), "https") {
		return fmt.Errorf("redirect URI must use HTTPS")
	}
	c := http.Client{}
	res, err := c.Get(fmt.Sprint(uri))
	if err != nil {
		return fmt.Errorf("one or more redirect URIs did not respond")
	}
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("one or more redirect URIs responded with a 404")
	}
	return nil
}
