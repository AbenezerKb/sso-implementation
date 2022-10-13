package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UserInfo struct {
	// Sub is unique and never reassigned identifier within for the End-User
	Sub string `json:"sub"`
	// FirstName is the first name of the user.
	FirstName string `json:"first_name,omitempty"`
	// MiddleName is the middle name of the user.
	MiddleName string `json:"middle_name,omitempty"`
	// LastName is the last name of the user.
	LastName string `json:"last_name,omitempty"`
	// Email is the email of the user.
	Email string `json:"email,omitempty"`
	// Phone is the phone of the user.
	Phone string `json:"phone,omitempty"`
	// Gender is the gender of the user.
	Gender string `json:"gender,omitempty"`
	// ProfilePicture is the profile image url for the user
	ProfilePicture string `json:"profile_picture,omitempty"`
}

func (u UserInfo) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required.Error("first name is required")),
		validation.Field(&u.Email, validation.When(u.Phone == "", validation.Required.Error("email is required"))),
		validation.Field(&u.Phone, validation.When(u.Email == "", validation.Required.Error("phone is required"))),
	)
}
