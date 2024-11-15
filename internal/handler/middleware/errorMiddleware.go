package middleware

import (
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"sso/internal/constant"
	"sso/internal/constant/errors"
	"sso/internal/constant/model"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/joomcode/errorx"
)

func ErrorHandler() gin.HandlerFunc {
	debugMode := viper.GetBool("debug")
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			e := c.Errors[0] // FIXME: how to handle multiple errors?
			err := e.Unwrap()
			for _, e := range errors.Error {
				if errorx.IsOfType(err, e.ErrorType) {
					er := errorx.Cast(err)
					response := model.ErrorResponse{
						Code:       e.ErrorCode,
						Message:    er.Message(),
						FieldError: ErrorFields(er.Cause()),
					}
					if debugMode {
						response.Description = fmt.Sprintf("Error: %v", er)
						response.StackTrace = fmt.Sprintf("%+v", errorx.EnsureStackTrace(err))
					}
					constant.ErrorResponse(c, &response)
					return
				}
			}
			constant.ErrorResponse(c, &model.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Unknown server error",
			})
			return
		}
	}
}
func ErrorFields(err error) []model.FieldError {
	var errs []model.FieldError
	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			errs = append(errs, model.FieldError{
				Name:        i,
				Description: v.Error(),
			},
			)
		}
		return errs
	}
	return nil
}
