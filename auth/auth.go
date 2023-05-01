package auth

import (
	"app/utils"
	"app/weberrors"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

var AdminToken = os.Getenv("ADMIN_TOKEN")

func Middleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		if token := gctx.Request.Header.Get(utils.API_AUTH_HEADER_KEY); token != AdminToken {
			gctx.AbortWithStatusJSON(http.StatusUnauthorized, weberrors.AppError{
				ErrorName:   http.StatusText(http.StatusUnauthorized),
				Description: "invalid admin token",
			})
			return
		}
		gctx.Next()
	}
}
