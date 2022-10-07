package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type ChangePhoneParam struct {
	// Phone is the phone of the user.
	Phone string `json:"phone"`
	// OTP is the one time password of the user.
	OTP string `json:"otp"`
}

func (c ChangePhoneParam) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Phone, validation.Required.Error("phone is required"), validation.By(validatePhone)),
		validation.Field(&c.OTP, validation.Required.Error("otp is required"), validation.Length(6, 6).Error("otp must be 6 characters")),
	)
}

type ChangePasswordParam struct {
	// NewPassword is the new password of the user.
	NewPassword string `json:"new_password,omitempty"`
	// OldPassword is the previous password of the user.
	OldPassword string `json:"old_password,omitempty"`
}

func (c ChangePasswordParam) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OldPassword, validation.Required.Error("old password is required")),
		validation.Field(&c.NewPassword, validation.Required.Error("new password is required"), validation.Length(6, 32).Error("password must be between 6 and 32 characters")),
	)
}
