package auth

import (
	"app/utils"
	"app/weberrors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	testFuncs "app/testing"
)

var AdminTokenTestString = "admin_token_string"

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
	originalToken := AdminToken
	AdminToken = AdminTokenTestString
	for _, testCase := range MiddlewareTestCases {
		t.Run(testCase.description, func(t *testing.T) {
			r := gin.New()
			r.Use(Middleware())
			r.Any("/", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, nil)
			})
			client := testFuncs.GetTestClient(t, r)
			res := client.GET("/").
				WithHeader(utils.API_AUTH_HEADER_KEY, testCase.adminToken).
				Expect()
			res.Status(testCase.expectedStatus)
			res.JSON().Equal(testCase.expectedResponse)
		})
	}
	AdminToken = originalToken
}
