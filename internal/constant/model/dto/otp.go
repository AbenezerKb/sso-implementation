package dto

type OTP struct {
	PhoneNumber string `json:"phone_number"`
	OTP         string `json:"otp"`
	Verified    bool   `json:"verified"`
	Type        string `json:"type"`
}