package weberrors

import (
	"net/http"
)

type AppError struct {
	ErrorName   string `json:"error"`
	Description string `json:"description"`
}

// func (a *AppError) Json() []byte {
// 	b, _ := json.Marshal(a)
// 	return b
// }

// func (a *AppError) String() string {
// 	j := a.Json()
// 	return string(j)
// }

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
const InternalServerDesc = "Internal Server Error."

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

var AvailableErrorReasons = map[string]string{
	"len":  "value is too long",
	"inv":  "value is invalid",
	"req":  "is required",
	"dup":  "already exists",
	"dupv": "has duplicate values",
	"miss": "value is missing",
	"inc":  "value is incorrect",
	"big":  "value is too big",
}
