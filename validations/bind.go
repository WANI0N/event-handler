package validations

import (
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var BindCustomValidators = func() {
	customValidations := []struct {
		FunctionTag string
		Function    validator.Func
	}{
		{"checkEventName", CheckEventNameValid},
		{"checkVideoQuality", CheckVideoQuality},
		{"checkAudioQuality", CheckAudioQuality},
		{"checkEmail", CheckEmailValid},
		{"checkTimeFieldFormat", CheckTimeFieldFormat},
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for _, validationDeclaration := range customValidations {
			err := v.RegisterValidation(
				validationDeclaration.FunctionTag,
				validationDeclaration.Function,
			)
			if err != nil {
				log.Logger.Error().
					Msgf("Failed to register custom validation `%v` - %v",
						validationDeclaration.FunctionTag,
						err,
					)
			}
		}
	}
}
