package dto

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
}
