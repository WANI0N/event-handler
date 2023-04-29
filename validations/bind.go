package validations

import (
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var BindCustomValidators = func() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("checkEventName", CheckEventNameValid)
		if err != nil {
			log.Logger.Error().
				Msgf("Failed to register checkEventName validation: %v", err)
		}
		err = v.RegisterValidation("checkVideoQuality", CheckVideoQuality)
		if err != nil {
			log.Logger.Error().
				Msgf("Failed to register checkVideoQuality validation: %v", err)
		}
		err = v.RegisterValidation("checkAudioQuality", CheckAudioQuality)
		if err != nil {
			log.Logger.Error().
				Msgf("Failed to register checkAudioQuality validation: %v", err)
		}
		err = v.RegisterValidation("checkEmail", CheckEmailValid)
		if err != nil {
			log.Logger.Error().
				Msgf("Failed to register checkEmail validation: %v", err)
		}
		err = v.RegisterValidation("checkTimeFieldFormat", CheckTimeFieldFormat)
		if err != nil {
			log.Logger.Error().
				Msgf("Failed to register checkTimeFieldFormat validation: %v", err)
		}
	}
}
