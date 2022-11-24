package request_models

import (
	"fmt"
	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type RSAPIUserRequest struct {
	// ID is the id of the user to be searched
	ID string `json:"id,omitempty" form:"id"`
	// Phone is the phone number of the user to be searched
	Phone string `json:"phone,omitempty" form:"phone"`
}

func (r RSAPIUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.When(r.Phone == "", validation.Required.Error("id is required"))),
		validation.Field(&r.Phone, validation.When(r.ID == "", validation.By(validatePhone))))
}

func validatePhone(phone interface{}) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}
