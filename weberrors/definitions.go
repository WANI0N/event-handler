package weberrors

import (
	"net/http"
)

type AppError struct {
	ErrorName   string `json:"error"`
	Description string `json:"description"`
}

type AppErrorWithCode struct {
	Code int `json:"code"`
	AppError
}

const ValidationErrorName = "ValidationError"
const InternalServerError = "InternalServerError"
const PayloadError = "PayloadError"
const NotFoundError = "NotFoundError"

const FieldErrorDescription = "Field `%v` %v"
const InvalidJsonPayloadDesc = "Invalid JSON payload."
const ResourceNotFoundErrorDesc = "The requested resource could not be found."
const RouteNotFoundErrorDesc = "Route does not exist."
const InternalServerDesc = "Internal Server Error."

var RouteNotFoundError = AppErrorWithCode{
	Code: http.StatusNotFound,
	AppError: AppError{
		ErrorName:   NotFoundError,
		Description: RouteNotFoundErrorDesc,
	},
}

var ValidationError = AppErrorWithCode{
	Code: http.StatusBadRequest,
	AppError: AppError{
		ErrorName:   ValidationErrorName,
		Description: FieldErrorDescription,
	},
}

var InvalidPayload = AppErrorWithCode{
	Code: http.StatusBadRequest,
	AppError: AppError{
		ErrorName:   PayloadError,
		Description: InvalidJsonPayloadDesc,
	},
}

var NotFound = AppErrorWithCode{
	Code: http.StatusNotFound,
	AppError: AppError{
		ErrorName:   NotFoundError,
		Description: ResourceNotFoundErrorDesc,
	},
}

var InternalError = AppErrorWithCode{
	Code: http.StatusInternalServerError,
	AppError: AppError{
		ErrorName:   InternalServerError,
		Description: InternalServerDesc,
	},
}
