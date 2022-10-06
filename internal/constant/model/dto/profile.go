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
