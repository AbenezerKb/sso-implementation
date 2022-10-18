package dto

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type Client struct {
	// ID is the unique identifier for the client.
	// It is automatically generated when the client is registered.
	ID uuid.UUID `json:"id"`
	// Name is the name of the client that will be displayed to the user.
	Name string `json:"name"`
	// ClientType is the type of the client.
	// It can be either confidential or public.
	ClientType string `json:"client_type"`
	// FirstParty shows if this client is a first party client or not
	FirstParty bool `json:"first_party"`
	// RedirectURIs is the list of redirect URIs of the client.
	// Each redirect URI must be a valid URL and must use HTTPS.
	RedirectURIs []string `json:"redirect_uris,omitempty"`
	// Scopes is the list of default scopes of the client if one is not provided.
	Scopes string `json:"scopes,omitempty"`
	// Secret is the secret the client uses to authenticate itself.
	// It is automatically generated when the client is registered.
	Secret string `json:"secret,omitempty"`
	// LogoURL is the URL of the client's logo.
	// It must be a valid URL.
	LogoURL string `json:"logo_url"`
	// Status is the current status of the client.
	// It is set to active by default.
	Status string `json:"status,omitempty"`
	// CreatedAt is the time this client was created at
	CreatedAt time.Time `json:"created_at"`
}

func (c Client) ValidateClient() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error("name is required"), validation.Length(3, 32).Error("name must be between 3 and 32 characters")),
		validation.Field(&c.ClientType, validation.Required.Error("client_type is required"), validation.In("confidential", "public").Error("client type must be either confidential or public")),
		validation.Field(&c.RedirectURIs, validation.Required.Error("redirect_uris is required"), validation.By(ValidateURI)),
		validation.Field(&c.Scopes, validation.Required.Error("scopes is required")),
		validation.Field(&c.LogoURL, validation.Required.Error("logo_url is required"), is.URL.Error("invalid logo_url"), validation.By(ValidateLogo)),
	)

}

func ValidateURI(uris interface{}) error {
	urisArray, ok := uris.([]string)
	if !ok {
		return fmt.Errorf("invalid uris")
	}
	for _, uri := range urisArray {
		if err := is.URL.Error("invalid redirect_uris").Validate(uri); err != nil {
			return err
		}

		if !strings.HasPrefix(fmt.Sprint(uri), "https") {
			return fmt.Errorf("redirect_uris must use https")
		}

		c := http.Client{}
		res, err := c.Get(fmt.Sprint(uri))
		if err != nil {
			return fmt.Errorf("redirect_uris not found")
		}
		if res.StatusCode == http.StatusNotFound {
			return fmt.Errorf("redirect_uris not found")
		}
	}

	return nil
}

func ValidateLogo(logo interface{}) error {
	c := http.Client{}
	res, err := c.Get(fmt.Sprint(logo))
	if err != nil || res.StatusCode != http.StatusOK {
		return fmt.Errorf("logo not found")
	}
	return nil
}

type UpdateClientStatus struct {
	// Status is new status that will replace old status of the user
	Status string `json:"status"`
}

func (u UpdateClientStatus) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Status, validation.Required.Error("status is required"), validation.In("ACTIVE", "PENDING", "INACTIVE")),
	)
}
