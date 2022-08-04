package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID             int64     `json:"id"`
	FirstName      string    `json:"first_name"`
	MiddleName     string    `json:"middle_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Password       string    `json:"password"`
	UserName       string    `json:"user_name"`
	Gender         string    `json:"gender"`
	Status         string    `json:"status"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

func (u User) ValidateUser() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required.Error("First name is required")),
		validation.Field(&u.LastName, validation.Required.Error("Last name is required")),
		validation.Field(&u.Email, is.Email.Error("Email is not valid")),
		validation.Field(&u.Phone, validation.Required.Error("Phone is required")),
		validation.Field(&u.Password, validation.Required.Error("Password is required"), validation.Length(6, 32).Error("Password must be between 6 and 32 characters")),
		validation.Field(&u.UserName, validation.Required.Error("User name is required")),
	)
}
