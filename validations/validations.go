package validations

import (
	"app/utils"
	"errors"
	"net/mail"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

var GetBindErrors = func(submitErrors error) error {
	var verrs validator.ValidationErrors
	if errors.As(submitErrors, &verrs) {
		return verrs
	}
	return nil
}

var allowedAudio = []string{"Low", "Mid", "High"}

var CheckAudioQuality validator.Func = func(fl validator.FieldLevel) bool {
	list := fl.Field().Interface().([]string)
	if len(list) == 0 {
		return true
	}
	for _, v := range list {
		if !slices.Contains(allowedAudio, v) {
			return false
		}
	}

	return true
}

var allowedResolutions = []string{"720p", "1080p", "1440p", "2160p"}

var CheckVideoQuality validator.Func = func(fl validator.FieldLevel) bool {
	list := fl.Field().Interface().([]string)
	if len(list) == 0 {
		return true
	}
	for _, v := range list {
		if !slices.Contains(allowedResolutions, v) {
			return false
		}
	}

	return true
}

var CheckUuidFormat = func(inputString string) bool {
	_, err := uuid.Parse(inputString)
	return err == nil
}

var CheckEventNameValid validator.Func = func(fl validator.FieldLevel) bool {
	return utils.EventNameRegex.MatchString(fl.Field().String())
}
var CheckEmailValid validator.Func = func(fl validator.FieldLevel) bool {
	emailList := fl.Field().Interface().([]string)
	var err error
	for _, email := range emailList {
		_, err = mail.ParseAddress(email)
		if err != nil {
			return false
		}
	}
	return true
}

const customTimeLayout = "2006-01-02T15:04:05Z"

var CheckTimeFieldFormat validator.Func = func(fl validator.FieldLevel) bool {
	_, err := time.Parse(customTimeLayout, fl.Field().String())
	return err == nil
}
