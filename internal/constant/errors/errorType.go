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
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrReadError,
	},
	{
		ErrorCode: http.StatusBadRequest,
		ErrorType: ErrDataExists,
	},
}

var (
	invalidInput = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	dbError      = errorx.NewNamespace("db error")
	duplicate    = errorx.NewNamespace("duplicate").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
)

var (
	ErrInvalidUserInput = errorx.NewType(invalidInput, "invalid user input")
	ErrNoRecordFound    = errorx.NewType(dbError, "no record found")
	ErrWriteError       = errorx.NewType(dbError, "could not write to db")
	ErrReadError        = errorx.NewType(dbError, "could not read from db")
	ErrDataExists       = errorx.NewType(duplicate, "data already exists")
)
