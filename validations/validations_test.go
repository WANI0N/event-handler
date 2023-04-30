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
	{"valid", "2023-04-20T14:00:00Z", true},
	{"invalid", "invalid-time-string", false},
}

var GetBindErrorsTestCases = []struct {
	description string
	haveError   bool
	matchResp   bool
}{
	{"get bind error", true, true},
	{"no bind error", false, false},
}

type inputStruct struct {
	a string
}

func TestGetBindErrors(t *testing.T) {
	for _, testCase := range GetBindErrorsTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validate := validator.New()
			validate.RegisterValidation("tag", func(fl validator.FieldLevel) bool {
				return testCase.haveError
			})
			err := validate.Var(inputStruct{a: ""}, "tag")

			resp := GetBindErrors(err)

			if testCase.matchResp {
				assert.Equal(t, resp, err)
			} else {
				assert.Nil(t, resp)
			}
		})
	}
}

var CheckVideoQualityTestCases = []struct {
	description  string
	submitList   []string
	expectedResp bool
}{
	{"Fail - invalid values", []string{"a", "b"}, false},
	{"Pass - empty list", []string{}, true},
	{"Pass - valid list", []string{"1080p", "1440p"}, true},
}

func TestCheckVideoQuality(t *testing.T) {
	for _, testCase := range CheckVideoQualityTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validate := validator.New()
			validate.RegisterValidation("checkVideoQuality", CheckVideoQuality)
			t.Run(testCase.description, func(t *testing.T) {
				err := validate.Var(testCase.submitList, "checkVideoQuality")
				if testCase.expectedResp {
					assert.Nil(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		})
	}
}

var CheckAudioQualityTestCases = []struct {
	description  string
	submitList   []string
	expectedResp bool
}{
	{"Fail - invalid values", []string{"a", "b"}, false},
	{"Pass - empty list", []string{}, true},
	{"Pass - valid list", []string{"Low", "High"}, true},
}

func TestCheckAudioQuality(t *testing.T) {
	for _, testCase := range CheckAudioQualityTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validate := validator.New()
			validate.RegisterValidation("checkAudioQuality", CheckAudioQuality)
			t.Run(testCase.description, func(t *testing.T) {
				err := validate.Var(testCase.submitList, "checkAudioQuality")
				if testCase.expectedResp {
					assert.Nil(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		})
	}
}

var CheckUuidFormatTestCases = []struct {
	description  string
	submitId     string
	expectedResp bool
}{
	{"valid", "0b1d34a7-8bcd-47f3-8923-d472510d8da4", true},
	{"invalid", "invalid-uuid-string", false},
}

func TestCheckUuidFormat(t *testing.T) {
	for _, testCase := range CheckUuidFormatTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			resp := CheckUuidFormat(testCase.submitId)
			assert.Equal(t, testCase.expectedResp, resp)
		})
	}
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
	{"Success", "valid event name", true},
	{"Fail", "invalid event name%?", false},
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
	{"1 valid email", []string{"validEmail@mail.com"}, true},
	{"multiple valid emails", []string{"a@mail.com", "b@mail.com"}, true},
	{"1st invalid email", []string{"invalidEmail"}, false},
	{"2nd invalid email", []string{"a@mail.com", "invalidEmail"}, false},
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
