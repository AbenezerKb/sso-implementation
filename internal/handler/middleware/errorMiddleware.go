package middleware

import (
	"fmt"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/joomcode/errorx"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			err := err.Unwrap()
			for _, e := range errors.Error {
				if errorx.IsOfType(err, e.ErrorType) {
					er := errorx.Cast(err)
					constant.ErrorResponse(c, &model.ErrorResponse{
						Code:        e.ErrorCode,
						Message:     er.Message(),
						FieldError:  Errorfields(er.Cause()),
						Description: fmt.Sprintf("Error: %v", er),
					})
				} else {
					constant.ErrorResponse(c, &model.ErrorResponse{
						Code:    http.StatusInternalServerError,
						Message: "Unknown server error",
					})
				}
			}
		}
	}
}
func Errorfields(err error) []string {
	var errors []string
	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			errors = append(errors, fmt.Sprintf("%v %v", i, v))
		}
		return errors
	}
	return nil
}
