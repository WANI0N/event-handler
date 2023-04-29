package auth

import (
	"app/utils"
	"app/weberrors"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
)

var MiddlewareTestCases = []struct {
	description      string
	adminToken       string
	expectedStatus   int
	expectedResponse interface{}
}{
	{
		description:    "Valid token",
		adminToken:     "admin_token_string",
		expectedStatus: http.StatusOK,
	},
	{
		description:    "Invalid token",
		adminToken:     "invalid_admin_token",
		expectedStatus: http.StatusUnauthorized,
		expectedResponse: &weberrors.AppError{
			ErrorName:   http.StatusText(http.StatusUnauthorized),
			Description: "invalid admin token",
		},
	},
}

func TestMiddleware(t *testing.T) {
	for _, testCase := range MiddlewareTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			r := gin.New()
			r.Use(Middleware())
			r.Any("/", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, nil)
			})
			client := httpexpect.WithConfig(httpexpect.Config{
				Client: &http.Client{
					Transport: httpexpect.NewBinder(r),
					Jar:       httpexpect.NewJar(),
				},
				Reporter: httpexpect.NewAssertReporter(t),
				Printers: []httpexpect.Printer{
					httpexpect.NewDebugPrinter(t, true),
				},
			})
			res := client.GET("/").
				WithHeader(utils.API_AUTH_HEADER_KEY, testCase.adminToken).
				Expect()
			res.Status(testCase.expectedStatus)
			res.JSON().Equal(testCase.expectedResponse)
		})
	}
}
