package validations

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var CheckTimeFieldFormatTestCases = []struct {
	description  string
	submitTime   string
	expectedResp bool
}{
	{
		description:  "Success",
		submitTime:   "2023-04-20T14:00:00Z",
		expectedResp: true,
	},
	{
		description:  "Fail",
		submitTime:   "invalid-time-string",
		expectedResp: false,
	},
}

func TestCheckTimeFieldFormat(t *testing.T) {
	for _, testCase := range CheckTimeFieldFormatTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validate := validator.New()
			validate.RegisterValidation("submitTime", CheckTimeFieldFormat)
			t.Run(testCase.description, func(t *testing.T) {
				err := validate.Var(testCase.submitTime, "submitTime")
				if testCase.expectedResp {
					assert.Nil(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		})
	}
}

var CheckEventNameValidTestCases = []struct {
	description  string
	submitName   string
	expectedResp bool
}{
	{
		description:  "Success",
		submitName:   "valid event name",
		expectedResp: true,
	},
	{
		description:  "Fail",
		submitName:   "invalid event name%?",
		expectedResp: false,
	},
}

func TestCheckEventNameValid(t *testing.T) {
	for _, testCase := range CheckEventNameValidTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validate := validator.New()
			validate.RegisterValidation("submitName", CheckEventNameValid)
			t.Run(testCase.description, func(t *testing.T) {
				err := validate.Var(testCase.submitName, "submitName")
				if testCase.expectedResp {
					assert.Nil(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		})
	}
}

var CheckEmailValidTestCases = []struct {
	description  string
	emailList    []string
	expectedResp bool
}{
	{
		description:  "Success, 1 valid email",
		emailList:    []string{"validEmail@mail.com"},
		expectedResp: true,
	},
	{
		description:  "Success, multiple valid emails",
		emailList:    []string{"validEmail@mail.com", "validEmail2@mail.com"},
		expectedResp: true,
	},
	{
		description:  "Fail, 1st invalid email",
		emailList:    []string{"invalidEmail"},
		expectedResp: false,
	},
	{
		description:  "Fail, 2nd invalid email",
		emailList:    []string{"validEmail@mail.com", "invalidEmail"},
		expectedResp: false,
	},
}

func TestCheckEmailValid(t *testing.T) {
	for _, testCase := range CheckEmailValidTestCases {
		validate := validator.New()
		validate.RegisterValidation("emailList", CheckEmailValid)
		t.Run(testCase.description, func(t *testing.T) {
			err := validate.Var(testCase.emailList, "emailList")
			if testCase.expectedResp {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
