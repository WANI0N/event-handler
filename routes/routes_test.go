package routes

import (
	"app/db"
	"app/models"
	"app/utils"
	"app/validations"
	"app/weberrors"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testClient(t *testing.T) *httpexpect.Expect {
	os.Setenv("CORS_ORIGIN", "*")
	app := gin.New()

	InitApp(app)
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(app),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	return e
}

func TestHealthCheckRoute(t *testing.T) {
	t.Run("Check if `ok` is returned in response", func(t *testing.T) {
		testClient := testClient(t)

		res := testClient.GET("/healthcheck").
			Expect()
		res.Status(http.StatusOK)
		res.Header("Content-type").Contains("application/json")
		res.JSON().Object().ValueEqual("result", "ok")
	})
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func getRandomAlphaNum(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

var CreateEventTestCases = []struct {
	description       string
	submitPayload     bool
	submitedPayload   interface{}
	dbCreateEventResp string
	dbCreateEventErr  error
	expectedStatus    int
	expectedResponse  interface{}
}{
	{
		description: "Success - default audio/video set",
		submitedPayload: models.EventData{
			Name:        "event-name",
			Timestamp:   "2023-04-20T14:00:00Z",
			Languages:   []string{"English"},
			Invitees:    []string{"valid-email@mail.com"},
			Description: "event-description",
		},
		dbCreateEventResp: "generated-uuid-string",
		expectedStatus:    http.StatusCreated,
		expectedResponse: models.EventResponseData{
			Id: "generated-uuid-string",
			EventData: models.EventData{
				Name:         "event-name",
				Timestamp:    "2023-04-20T14:00:00Z",
				VideoQuality: []string{utils.DEFAULT_RESOLUTION},
				AudioQuality: []string{utils.DEVAULT_AUDIO},
				Languages:    []string{"English"},
				Invitees:     []string{"valid-email@mail.com"},
				Description:  "event-description",
			},
		},
	},
	{
		description: "Success - user set audio/video params",
		submitedPayload: models.EventData{
			Name:         "event-name",
			Timestamp:    "2023-04-20T14:00:00Z",
			Languages:    []string{"English"},
			VideoQuality: []string{"2160p", "1440p"},
			AudioQuality: []string{"High", "Mid"},
			Invitees:     []string{"valid-email@mail.com", "valid-email2@mail.com"},
			Description:  "event-description",
		},
		dbCreateEventResp: "generated-uuid-string",
		expectedStatus:    http.StatusCreated,
		expectedResponse: models.EventResponseData{
			Id: "generated-uuid-string",
			EventData: models.EventData{
				Name:         "event-name",
				Timestamp:    "2023-04-20T14:00:00Z",
				VideoQuality: []string{"2160p", "1440p"},
				AudioQuality: []string{"High", "Mid"},
				Languages:    []string{"English"},
				Invitees:     []string{"valid-email@mail.com", "valid-email2@mail.com"},
				Description:  "event-description",
			},
		},
	},
	{
		description:    "Fail - required fields validation",
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			"Field `name` is required, field `timestamp` is required, field `languages` is required, field `invitees` is required.")),
	},
	{
		description:      "Fail - invalid Json payload",
		submitedPayload:  "not a json string",
		expectedStatus:   http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(&weberrors.InvalidPayload),
	},
	{
		description: "Fail - incorrect `name` format",
		submitedPayload: models.EventData{
			Name:      "event-name????",
			Timestamp: "2023-04-20T14:00:00Z",
			Languages: []string{"English"},
			Invitees:  []string{"valid-email@mail.com"},
		},
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			"field `name` contains invalid characters (use A-Za-z0-9 _- only)")),
	},
	{
		description: "Fail - incorrect `email` format",
		submitedPayload: models.EventData{
			Name:      "event-name",
			Timestamp: "2023-04-20T14:00:00Z",
			Languages: []string{"English"},
			Invitees:  []string{"invalid-email"},
		},
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			"field `invitees` contains invalid email address")),
	},
	{
		description: "Fail - incorrect `date` format",
		submitedPayload: models.EventData{
			Name:      "event-name",
			Timestamp: "2023-04-20 ??? 14:00:00Z",
			Languages: []string{"English"},
			Invitees:  []string{"valid-email@mail.com"},
		},
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			"field date does not have correct format (use YYYY-MM-DDTHH:MM:SSZ)")),
	},
	{
		description: "Fail - incorrect `videoQuality` format",
		submitedPayload: models.EventData{
			Name:         "event-name",
			Timestamp:    "2023-04-20T14:00:00Z",
			VideoQuality: []string{"123p"},
			Languages:    []string{"English"},
			Invitees:     []string{"valid-email@mail.com"},
		},
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			fmt.Sprintf(
				"field `videoQuality` contains invalid resolution (allowed values: %v)",
				strings.Join(utils.ALLOWED_RESOLUTION, ", "),
			))),
	},
	{
		description: "Fail - incorrect `audioQuality` format",
		submitedPayload: models.EventData{
			Name:         "event-name",
			Timestamp:    "2023-04-20T14:00:00Z",
			AudioQuality: []string{"very-nice"},
			Languages:    []string{"English"},
			Invitees:     []string{"valid-email@mail.com"},
		},
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			fmt.Sprintf(
				"field `audioQuality` contains invalid resolution (allowed values: %v)",
				strings.Join(utils.ALLOWED_AUDIO, ", "),
			))),
	},
	{
		description: "Fail - field `name` too long; field `invitees` has duplicates",
		submitedPayload: models.EventData{
			Name:      getRandomAlphaNum(256),
			Timestamp: "2023-04-20T14:00:00Z",
			Languages: []string{"English"},
			Invitees:  []string{"valid-email@mail.com", "valid-email@mail.com"},
		},
		expectedStatus: http.StatusBadRequest,
		expectedResponse: weberrors.ParseAppError(weberrors.ValidationError.ChangeDesc(
			"field `name` cannot be longer than 255, field `invitees` contains duplicate values")),
	},
}

func TestCreateEventRoute(t *testing.T) {
	for _, testCase := range CreateEventTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			db.CreateEvent = func(payload models.EventData) (string, error) {
				convertedTestCaseData := testCase.submitedPayload.(models.EventData)
				if len(convertedTestCaseData.VideoQuality) == 0 {
					convertedTestCaseData.VideoQuality = []string{utils.DEFAULT_RESOLUTION}
				}
				if len(convertedTestCaseData.AudioQuality) == 0 {
					convertedTestCaseData.AudioQuality = []string{utils.DEVAULT_AUDIO}
				}
				assert.Equal(t, convertedTestCaseData, payload)
				return testCase.dbCreateEventResp, testCase.dbCreateEventErr
			}

			res := testClient(t).POST("/event").
				WithJSON(testCase.submitedPayload).Expect()
			res.Header("Content-type").Contains("application/json")
			res.Status(testCase.expectedStatus)
			res.JSON().Equal(testCase.expectedResponse)
		})
	}
}

var GetEventTestCases = []struct {
	description                    string
	submitIdPathParam              string
	validationsCheckUuidFormatResp bool
	dbGetEventMockResp             models.EventResponseData
	dbGetEventMockErr              error
	expectedStatus                 int
	expectedResp                   interface{}
}{
	{
		description:                    "Success",
		submitIdPathParam:              "90a04b08-d820-4106-8ced-2cbc940728a3",
		validationsCheckUuidFormatResp: true,
		dbGetEventMockResp: models.EventResponseData{
			Id: "90a04b08-d820-4106-8ced-2cbc940728a3",
			EventData: models.EventData{
				Name: "event-name",
			},
		},
		expectedStatus: http.StatusOK,
		expectedResp: models.EventResponseData{
			Id: "90a04b08-d820-4106-8ced-2cbc940728a3",
			EventData: models.EventData{
				Name: "event-name",
			},
		},
	},
	{
		description:                    "Fail - invalid uuid - resource not found",
		submitIdPathParam:              "invalid-uuid-string",
		validationsCheckUuidFormatResp: false,
		expectedStatus:                 http.StatusNotFound,
		expectedResp:                   weberrors.ParseAppError(&weberrors.NotFound),
	},
	{
		description:                    "Fail - event does not exist",
		submitIdPathParam:              "90a04b08-d820-4106-8ced-2cbc940728a3",
		validationsCheckUuidFormatResp: true,
		dbGetEventMockErr:              errors.New("not found"),
		expectedStatus:                 http.StatusNotFound,
		expectedResp:                   weberrors.ParseAppError(&weberrors.NotFound),
	},
	{
		description:                    "Fail - db unexpected error",
		submitIdPathParam:              "90a04b08-d820-4106-8ced-2cbc940728a3",
		validationsCheckUuidFormatResp: true,
		dbGetEventMockErr:              errors.New("redis connection error"),
		expectedStatus:                 http.StatusInternalServerError,
		expectedResp:                   weberrors.ParseAppError(&weberrors.InternalError),
	},
}

func TestGetEvent(t *testing.T) {
	for _, testCase := range GetEventTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validations.CheckUuidFormat = func(inputString string) bool {
				assert.Equal(t, testCase.submitIdPathParam, inputString)
				return testCase.validationsCheckUuidFormatResp
			}
			db.GetEvent = func(id string) (models.EventResponseData, error) {
				assert.Equal(t, testCase.submitIdPathParam, id)
				return testCase.dbGetEventMockResp, testCase.dbGetEventMockErr
			}
			res := testClient(t).GET(
				fmt.Sprintf("/event/%v", testCase.submitIdPathParam)).Expect()
			res.Header("Content-type").Contains("application/json")
			res.Status(testCase.expectedStatus)
			res.JSON().Equal(testCase.expectedResp)
		})
	}
}

var DeleteEventTestCases = []struct {
	description                    string
	submitIdPathParam              string
	validationsCheckUuidFormatResp bool
	adminToken                     string
	dbDeleteEventMockErr           error
	expectedStatus                 int
	expectedResp                   interface{}
}{
	{
		description:                    "Success",
		submitIdPathParam:              "90a04b08-d820-4106-8ced-2cbc940728a3",
		adminToken:                     utils.GetEnvOrDefault("ADMIN_TOKEN", "admin_token_string"),
		validationsCheckUuidFormatResp: true,
		expectedStatus:                 http.StatusNoContent,
	},
	{
		description:                    "Fail - Unauthorized",
		submitIdPathParam:              "1",
		adminToken:                     "invalid_admin_token",
		validationsCheckUuidFormatResp: true,
		expectedStatus:                 http.StatusUnauthorized,
		expectedResp: &weberrors.AppError{
			ErrorName:   http.StatusText(http.StatusUnauthorized),
			Description: "invalid admin token",
		},
	},
	{
		description:                    "Fail - db unexpected error",
		submitIdPathParam:              "90a04b08-d820-4106-8ced-2cbc940728a3",
		adminToken:                     utils.GetEnvOrDefault("ADMIN_TOKEN", "admin_token_string"),
		validationsCheckUuidFormatResp: true,
		dbDeleteEventMockErr:           errors.New("redis connection error"),
		expectedStatus:                 http.StatusInternalServerError,
		expectedResp:                   weberrors.ParseAppError(&weberrors.InternalError),
	},
}

func TestDeleteEvent(t *testing.T) {
	for _, testCase := range DeleteEventTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			validations.CheckUuidFormat = func(inputString string) bool {
				assert.Equal(t, testCase.submitIdPathParam, inputString)
				return testCase.validationsCheckUuidFormatResp
			}
			db.DeleteEvent = func(id string) error {
				assert.Equal(t, testCase.submitIdPathParam, id)
				return testCase.dbDeleteEventMockErr
			}
			res := testClient(t).DELETE(
				fmt.Sprintf("/event/%v", testCase.submitIdPathParam)).
				WithHeader(utils.API_AUTH_HEADER_KEY, testCase.adminToken).
				Expect()

			res.Status(testCase.expectedStatus)
			if testCase.expectedResp != nil {
				res.JSON().Equal(testCase.expectedResp)
			}
		})
	}
}
