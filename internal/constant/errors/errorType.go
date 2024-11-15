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
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrDBDelError,
	},
	{
		ErrorCode: http.StatusBadRequest,
		ErrorType: ErrDataExists,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrCacheSetError,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrCacheGetError,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrCacheDel,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrInternalServerError,
	},
	{
		ErrorCode: http.StatusUnauthorized,
		ErrorType: ErrInvalidToken,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrOTPGenerate,
	},
	{
		ErrorCode: http.StatusInternalServerError,
		ErrorType: ErrSMSSend,
	},
	{
		ErrorCode: http.StatusUnauthorized,
		ErrorType: ErrAuthError,
	},
	{
		ErrorCode: http.StatusForbidden,
		ErrorType: ErrAcessError,
	},
}

var (
	invalidInput = errorx.NewNamespace("validation error").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	dataNotFound = errorx.NewNamespace("data not found").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	unauthorized = errorx.NewNamespace("unauthorized").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	dbError      = errorx.NewNamespace("db error")
	duplicate    = errorx.NewNamespace("duplicate").ApplyModifiers(errorx.TypeModifierOmitStackTrace)
	cacheError   = errorx.NewNamespace("cache error")
	serverError  = errorx.NewNamespace("server error")
	AccessDenied = errorx.RegisterTrait("You are not authorized to perform the action")
	kafkaError   = errorx.NewNamespace("kafka error")
)

var (
	ErrInvalidUserInput    = errorx.NewType(invalidInput, "invalid user input")
	ErrNoRecordFound       = errorx.NewType(dataNotFound, "no record found")
	ErrWriteError          = errorx.NewType(dbError, "could not write to db")
	ErrReadError           = errorx.NewType(dbError, "could not read from db")
	ErrDBDelError          = errorx.NewType(dbError, "could not delete record")
	ErrUpdateError         = errorx.NewType(dbError, "could not update db")
	ErrDataExists          = errorx.NewType(duplicate, "data already exists")
	ErrCacheSetError       = errorx.NewType(cacheError, "could not set cache")
	ErrCacheGetError       = errorx.NewType(cacheError, "could not get cache")
	ErrCacheDel            = errorx.NewType(cacheError, "could not delete cache")
	ErrInternalServerError = errorx.NewType(serverError, "internal server error")
	ErrInvalidToken        = errorx.NewType(unauthorized, "invalid token")
	ErrOTPGenerate         = errorx.NewType(serverError, "couldn't generate otp")
	ErrSMSSend             = errorx.NewType(serverError, "couldn't send sms")
	ErrAuthError           = errorx.NewType(unauthorized, "you are not authorized.")
	ErrAcessError          = errorx.NewType(errorx.CommonErrors, "Unauthorized", AccessDenied)
	ErrKafkaRead           = errorx.NewType(kafkaError, "could not read from kafka")
	ErrKafkaInvalidEvent   = errorx.NewType(kafkaError, "invalid kafka event")
)
