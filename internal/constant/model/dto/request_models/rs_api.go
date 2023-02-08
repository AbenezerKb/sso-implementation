package request_models

import (
	"fmt"

	"github.com/dongri/phonenumber"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RSAPIUserRequest struct {
	// ID is the id of the user to be searched
	ID string `json:"id,omitempty" form:"id"`
	// Phone is the phone number of the user to be searched
	Phone string `json:"phone,omitempty" form:"phone"`
}

type RSAPIUsersRequest struct {
	// IDs is the ids of the users to be searched
	IDs []string `json:"ids,omitempty"`
	// Phones is the phone numbers of the users to be searched
	Phones []string `json:"phones,omitempty"`
}

func (r RSAPIUsersRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.IDs,
			//validation.When(len(r.Phones) == 0, validation.Required.Error("ids or phones is required")),
			validation.When(len(r.IDs) > 0, validation.Each(is.UUID))),
		validation.Field(&r.Phones,
			//validation.When(len(r.IDs) == 0, validation.Required.Error("ids or phones is required")),
			validation.When(len(r.Phones) > 0, validation.Each(validation.By(validatePhone)))),
	)
}

func (r RSAPIUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.When(r.Phone == "", validation.Required.Error("id or phone is required"))),
		validation.Field(&r.Phone, validation.When(r.ID == "", validation.By(validatePhone))))
}

func validatePhone(phone interface{}) error {
	str := phonenumber.Parse(fmt.Sprintf("%v", phone), "ET")
	if str == "" {
		return fmt.Errorf("invalid phone number")
	}
	return nil
}
