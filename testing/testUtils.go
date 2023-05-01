package testing

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/gin-gonic/gin"
)

func GetTestClient(testSuite *testing.T, appEngine *gin.Engine) *httpexpect.Expect {
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(appEngine),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(testSuite),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(testSuite, true),
		},
	})
	return e
}
