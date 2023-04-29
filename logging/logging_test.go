package logging

import (
	"app/utils"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func captureLogOutput(f func()) string {
	var buf bytes.Buffer
	log.Logger = log.Output(&buf)
	f()
	log.Logger = log.Output(os.Stderr)
	return buf.String()
}

func TestSetupDefaultLogger(t *testing.T) {

	t.Run("Should setup default field names", func(t *testing.T) {
		assert.Equal(t, timestampFieldName, zerolog.TimestampFieldName,
			"should set timestamp field name")
	})

	t.Run("Should add fields to default logger", func(t *testing.T) {
		output := captureLogOutput(func() {
			log.Info().Msg("Testing")
		})

		assert.Contains(t, output, fmt.Sprintf("%q:", timestampFieldName),
			"should add timestamp")
		assert.Contains(t, output, fmt.Sprintf("%q:", "caller"),
			"should add caller")
		assert.Contains(t, output, fmt.Sprintf("%q:%q", "service", utils.APP_NAME),
			"should add service name")
		assert.Contains(t, output, fmt.Sprintf("%q:", "environment"),
			"should add environment name")
	})
}

var ginLoggerMiddlewareTestCases = []struct {
	description             string
	responseStatus          int
	expectedHighestLogLevel zerolog.Level
}{
	{
		description:             "log with info level",
		responseStatus:          http.StatusOK,
		expectedHighestLogLevel: zerolog.InfoLevel,
	},
	{
		description:             "log with warn level if client error",
		responseStatus:          http.StatusBadRequest,
		expectedHighestLogLevel: zerolog.WarnLevel,
	},
	{
		description:             "log with error level if server error",
		responseStatus:          http.StatusInternalServerError,
		expectedHighestLogLevel: zerolog.ErrorLevel,
	},
}

func TestMiddleware(t *testing.T) {
	for _, testCase := range ginLoggerMiddlewareTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
			request, _ := http.NewRequest(http.MethodGet, "/", bytes.NewBuffer([]byte("{}")))
			ginContext.Request = request
			ginContext.Status(testCase.responseStatus)
			middlewareHandler := Middleware()

			logOutput := captureLogOutput(func() {
				middlewareHandler(ginContext)
			})

			correlationId, _ := ginContext.Get(correlationIdContextKey)
			assert.NotEmpty(t, correlationId)

			assert.Contains(t, logOutput, requestStartedMessage,
				"should log start of the request")
			assert.Contains(t, logOutput, correlationId,
				"should add correlation id to logs")
			assert.Contains(t, logOutput, request.URL.Path,
				"should add called url to logs")
			assert.Contains(t, logOutput, request.URL.RawQuery,
				"should include query")
			assert.Contains(t, logOutput, requestEndedMessage,
				"should log end of request")
			assert.Contains(t, logOutput, testCase.expectedHighestLogLevel.String(),
				"should log outbound log entry with correct level")

		})
	}
}
