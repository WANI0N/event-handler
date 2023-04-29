package weberrors

import (
	"app/logging"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func JSONAppErrorReporter() gin.HandlerFunc {
	return customJsonAppErrorReporter(gin.ErrorTypeAny)
}

func ParseAppError(err error) AppError {
	if err == nil {
		return AppError{}
	}
	e := strings.SplitN(err.Error(), "-", 2)
	if len(e) < 2 {
		panic("`error` has to `AppError` type")
	}
	appError := AppError{
		ErrorName:   e[0],
		Description: e[1],
	}
	return appError
}

func customJsonAppErrorReporter(errType gin.ErrorType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		detectedErrors := ctx.Errors.ByType(errType)

		if len(detectedErrors) == 0 {
			return
		}

		logger := logging.WithContext(ctx)
		err := detectedErrors[0].Err
		var parsedError *AppError
		var ErrorName int = 500

		switch err := err.(type) {
		case validator.ValidationErrors:
			logger.Error().Err(err).Msg("Request validation error")
			ctx.JSON(
				http.StatusBadRequest,
				AppError{
					ErrorName:   ValidationErrorName,
					Description: GetErrorText(err),
				})
			return
		case *AppErrorWithCode:
			ErrorName = err.Code
			parsedError = &AppError{
				ErrorName:   err.ErrorName,
				Description: err.Description,
			}
		default:
			logger.Error().Err(err).Msg("Unexpected error occurred")
			parsedError = &AppError{
				ErrorName:   InternalServerError,
				Description: "Internal Server Error.",
			}
		}
		ctx.JSON(ErrorName, parsedError)
	}
}
