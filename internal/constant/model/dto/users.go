package dto

import (
	"github.com/google/uuid"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID             uuid.UUID `json:"id"`
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
		validation.Field(&u.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&u.MiddleName, validation.Required.Error("middle name is required")),
		validation.Field(&u.LastName, validation.Required.Error("last name is required")),
		validation.Field(&u.Email, is.Email.Error("email is not valid")),
		validation.Field(&u.Phone, validation.Required.Error("phone is required")),
		validation.Field(&u.Password, validation.Required.Error("password is required"), validation.Length(6, 32).Error("password must be between 6 and 32 characters")),
		validation.Field(&u.UserName, validation.Required.Error("user name is required")),
	)
}
