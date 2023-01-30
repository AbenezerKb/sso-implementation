package dto

import (
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UploadAssetRequest struct {
	Asset *multipart.FileHeader `form:"asset" json:"asset"`
	Type  string                `form:"type" json:"type"`
}

func (u UploadAssetRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Type, validation.Required.Error("type is required")),
		validation.Field(&u.Asset, validation.Required.Error("asset is required")))
}
