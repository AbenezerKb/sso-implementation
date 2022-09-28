// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AuthHistory struct {
	ID          uuid.UUID      `json:"id"`
	Code        string         `json:"code"`
	UserID      uuid.UUID      `json:"user_id"`
	Scope       sql.NullString `json:"scope"`
	Status      string         `json:"status"`
	RedirectUri sql.NullString `json:"redirect_uri"`
	ClientID    uuid.UUID      `json:"client_id"`
	CreatedAt   time.Time      `json:"created_at"`
}

type Client struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	ClientType   string    `json:"client_type"`
	RedirectUris string    `json:"redirect_uris"`
	Scopes       string    `json:"scopes"`
	Secret       string    `json:"secret"`
	LogoUrl      string    `json:"logo_url"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type Internalrefreshtoken struct {
	ID           uuid.UUID `json:"id"`
	Refreshtoken string    `json:"refreshtoken"`
	UserID       uuid.UUID `json:"user_id"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RefreshToken struct {
	ID           uuid.UUID      `json:"id"`
	RefreshToken string         `json:"refresh_token"`
	Code         string         `json:"code"`
	UserID       uuid.UUID      `json:"user_id"`
	Scope        sql.NullString `json:"scope"`
	RedirectUri  sql.NullString `json:"redirect_uri"`
	ExpiresAt    time.Time      `json:"expires_at"`
	ClientID     uuid.UUID      `json:"client_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type ResourceServer struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Scope struct {
	ID                 uuid.UUID      `json:"id"`
	Name               string         `json:"name"`
	Description        string         `json:"description"`
	ResourceServerID   uuid.NullUUID  `json:"resource_server_id"`
	ResourceServerName sql.NullString `json:"resource_server_name"`
	Status             string         `json:"status"`
}

type User struct {
	ID             uuid.UUID      `json:"id"`
	FirstName      string         `json:"first_name"`
	MiddleName     string         `json:"middle_name"`
	LastName       string         `json:"last_name"`
	Email          sql.NullString `json:"email"`
	Phone          string         `json:"phone"`
	Password       string         `json:"password"`
	UserName       string         `json:"user_name"`
	Gender         string         `json:"gender"`
	ProfilePicture sql.NullString `json:"profile_picture"`
	Status         sql.NullString `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
}
