package dto

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID             uuid.UUID `json:"id,omitempty"`
	FirstName      string    `json:"first_name,omitempty"`
	MiddleName     string    `json:"middle_name,omitempty"`
	LastName       string    `json:"last_name,omitempty"`
	Email          string    `json:"email,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	Password       string    `json:"password,omitempty"`
	UserName       string    `json:"user_name,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	Status         string    `json:"status,omitempty"`
	ProfilePicture string    `json:"profile_picture,omitempty"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
	OTP            string    `json:"otp,omitempty"`
}

func (u User) ValidateUser() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&u.MiddleName, validation.Required.Error("middle name is required")),
		validation.Field(&u.LastName, validation.Required.Error("last name is required")),
		validation.Field(&u.Email, is.EmailFormat.Error("email is not valid")),
		validation.Field(&u.Phone, validation.Required.Error("phone is required"), validation.By(validatePhone)),
		validation.Field(&u.Password, validation.When(u.Email != "", validation.Required.Error("password is required"), validation.Length(6, 32).Error("password must be between 6 and 32 characters"))),
		validation.Field(&u.OTP, validation.Required.Error("otp is required"), validation.Length(6, 6).Error("otp must be 6 characters")),
	)
}
func (u User) ValidateLoginCredentials() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Phone, validation.When(u.OTP != "" && u.Email == "",
			validation.Required.Error("phone is required"),
			validation.By(validatePhone))),
		validation.Field(&u.OTP, validation.When(u.Phone != "",
			validation.Required.Error("otp is required"),
			validation.Length(6, 6).Error("otp must be 6 characters"))),
		validation.Field(&u.Email, validation.When(u.Phone == "" && u.Password != "",
			validation.Required.Error("email is required"),
			is.EmailFormat.Error("email is not valid"))),
		validation.Field(&u.Password, validation.When(u.Email != "",
			validation.Required.Error("password is required"),
			validation.Length(6, 32).Error("password must be between 6 and 32 characters"))),
	)
}
func validatePhone(phone interface{}) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}
