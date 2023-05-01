package weberrors

import (
	"app/utils"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (e *AppErrorWithCode) Error() string {
	return fmt.Sprintf("%v-%v", e.ErrorName, e.Description)
}

func (e *AppErrorWithCode) ChangeDesc(newDescription string) error {
	var err = *e

	wordArray := strings.Split(strings.TrimSpace(newDescription), " ")
	wordArray[0] = cases.Title(language.English).String(wordArray[0])
	newDescription = strings.Join(wordArray[:], " ")

	if string(newDescription[len(newDescription)-1]) != "." {
		newDescription += "."
	}

	err.AppError.Description = newDescription
	return &err
}

var ValidationErrorToText = func(e validator.FieldError) string {
	field := e.Field()
	field = strings.ToLower(field[0:1]) + field[1:]
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("field `%s` is required", field)
	case "len":
		return fmt.Sprintf("field `%s` must be of length %s", field, e.Param())
	case "max":
		return fmt.Sprintf("field `%s` cannot be longer than %s", field, e.Param())
	case "min":
		return fmt.Sprintf("field `%s` must be longer than %s", field, e.Param())
	case "unique":
		return fmt.Sprintf("field `%s` contains duplicate values", field)
	case "oneof":
		return fmt.Sprintf("field `%s` needs to be one of values: %s", field, e.Param())
	case "checkEmail":
		return fmt.Sprintf("field `%s` contains invalid email address", field)
	case "checkVideoQuality":
		return fmt.Sprintf(
			"field `%s` contains invalid resolution (allowed values: %v)",
			field, strings.Join(utils.ALLOWED_RESOLUTION, ", "),
		)
	case "checkAudioQuality":
		return fmt.Sprintf("field `%s` contains invalid resolution (allowed values: %v)",
			field, strings.Join(utils.ALLOWED_AUDIO, ", "),
		)
	case "checkEventName":
		return fmt.Sprintf(
			"field `%s` contains invalid characters (use A-Za-z0-9 _- only)", field)
	case "checkTimeFieldFormat":
		return fmt.Sprint(
			"field date does not have correct format (use YYYY-MM-DDTHH:MM:SSZ)")
	}

	return fmt.Sprintf("field `%s` is invalid", field)
}

func GetErrorText(verrs validator.ValidationErrors) string {
	l := []string{}
	for i, e := range verrs {
		msg := ValidationErrorToText(e)

		if i == 0 {
			msg = Capitalize(msg)
		}
		l = append(l, msg)
	}

	if len(l) == 0 {
		return ""
	}

	return strings.Join(l, ", ") + "."
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return cases.Title(language.Und).String(s[0:1]) + s[1:]
}
