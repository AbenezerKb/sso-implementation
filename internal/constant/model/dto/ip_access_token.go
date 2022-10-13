package dto

import (
	"github.com/google/uuid"
	"time"
)

// IPAccessToken holds access token (refresh_token too) for a login with an identity provider
type IPAccessToken struct {
	// ID is the unique id for this access token
	ID uuid.UUID `json:"id"`
	// UserID is the id of the user this access token belongs to
	UserID uuid.UUID `json:"user_id"`
	// SubID is the unique identifier for the user from the identity provider
	SubID string `json:"sub_id"`
	// IPID is the id of the identity provider this access token is granted by
	IPID uuid.UUID `json:"ip_id"`
	// Token is the actual access token string
	Token string `json:"token"`
	// RefreshToken is an optional refresh token to get a new access token with
	RefreshToken string `json:"refresh_token,omitempty"`
	// Status is the current status of this access token
	Status string `json:"status"`
	// CreatedAt is the time this access token is first created at
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt is the time this access token is last updated at
	UpdatedAt time.Time `json:"updated_at"`
}
