package dto

import (
	"github.com/google/uuid"
	"time"
)

// ResourceServer is a server that this sso controls access for
type ResourceServer struct {
	// ID is the unique id for this resource server
	ID uuid.UUID `json:"id"`
	// Name is the resource server name.
	// It must be unique across the sso
	Name string `json:"name"`
	// CreatedAt is the time this resource server is created at
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time this resource server is updated at
	UpdatedAt time.Time `json:"updated_at"`
}
