package weberrors

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var ChangeDescReturnsNewDescTestCases = []struct {
	Description  string
	SubmitDesc   string
	ExpectedDesc string
}{
	{
		Description:  "No alterations made",
		SubmitDesc:   "Capital and period present.",
		ExpectedDesc: "Capital and period present.",
	},
	{
		Description:  "Alterations made",
		SubmitDesc:   "no capital and no period",
		ExpectedDesc: "No capital and no period.",
	},
}

func TestChangeDesc(t *testing.T) {
	for _, testCase := range ChangeDescReturnsNewDescTestCases {
		t.Run(testCase.Description, func(t *testing.T) {
			exp_err := fmt.Errorf("%v-%v", ValidationErrorName, testCase.ExpectedDesc)
			err := ValidationError.ChangeDesc(testCase.SubmitDesc)
			assert.Equal(t, exp_err.Error(), err.Error())
		})
	}
}

func TestGetErrorText(t *testing.T) {

	t.Run("without errors, should return empty message", func(t *testing.T) {
		input := []validator.FieldError{}

		r := GetErrorText(input)

		assert.Equal(t, "", r)
	})

	t.Run("with single error, should return formatted message", func(t *testing.T) {
		type TestSingleError struct {
			Required int `validate:"required"`
		}
		validate := validator.New()
		errs := validate.Struct(TestSingleError{})
		assert.NotNil(t, errs)
		verrs := errs.(validator.ValidationErrors)
		assert.NotEmpty(t, verrs)

		r := GetErrorText(verrs)

		assert.Equal(t, "Field `required` is required.", r)
	})

	t.Run("with multiple errors, should return formatted message", func(t *testing.T) {
		type TestSingleError struct {
			Required int `validate:"required"`
			Name     int `validate:"min=1"`
		}
		validate := validator.New()
		errs := validate.Struct(TestSingleError{})
		assert.NotNil(t, errs)
		verrs := errs.(validator.ValidationErrors)
		assert.NotEmpty(t, verrs)

		r := GetErrorText(verrs)

		assert.Equal(t, "Field `required` is required, field `name` must be longer than 1.", r)
	})
}

func TestCapitalize(t *testing.T) {
	t.Run("with whitespace, should pass through", func(t *testing.T) {
		assert.Equal(t, "", Capitalize(""))
		assert.Equal(t, " ", Capitalize(" "))
		assert.Equal(t, " a", Capitalize(" a"))
	})

	t.Run("with a text, should capitalize", func(t *testing.T) {
		assert.Equal(t, "Lorem ipsum Dolor sit Amet", Capitalize("lorem ipsum Dolor sit Amet"))
	})
}
