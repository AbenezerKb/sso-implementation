package request_models

import (
	"time"

	"github.com/google/uuid"
)

type MinRideEvent struct {
	Event  string
	Driver *Driver
}

type MiniRideDriverResponse struct {
	ID             uuid.UUID `json:"id,omitempty"`
	FullName       string    `json:"full_name,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	Status         string    `json:"status,omitempty"`
	DriverID       uuid.UUID `json:"driver_id"`
	DriverLicense  string    `json:"driver_license"`
	State          string    `json:"state"`
	SwapPhones     []string  `json:"swap_phones,omitempty"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type Driver struct {
	ID             uuid.UUID `json:"id,omitempty"`
	FirstName      string    `json:"first_name,omitempty"`
	MiddleName     string    `json:"middle_name,omitempty"`
	LastName       string    `json:"last_name,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	Status         string    `json:"status,omitempty"`
	DriverID       uuid.UUID `json:"driverId"`

	SwapPhones []string  `json:"swap_phones,omitempty"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}
