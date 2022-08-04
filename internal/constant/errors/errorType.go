package errors

import (
	"net/http"

	"github.com/joomcode/errorx"
)

type ErrorType struct {
	ErrorCode int
	ErrorType *errorx.Type
}

var Error = []ErrorType{
	{
		ErrorCode: http.StatusBadRequest,
		ErrorType: ErrInvalidUserInput,
	},
	{
		ErrorCode: http.StatusNotFound,
		ErrorType: ErrNoRecordFound,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrWriteError,
	},
}

var (
	invalidInput  = errorx.NewNamespace("validation error")
	dbError = errorx.NewNamespace("db error")
)

var (
	ErrInvalidUserInput = errorx.NewType(invalidInput, "invalid user input")
	ErrNoRecordFound      = errorx.NewType(dbError, "no record found")
	ErrWriteError         = errorx.NewType(dbError, "could not write to db")
)
