package auth

import (
	"app/utils"
	"app/weberrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var adminToken = utils.GetEnvOrDefault("ADMIN_TOKEN", "admin_token_string")

func Middleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		if token := gctx.Request.Header.Get(utils.API_AUTH_HEADER_KEY); token != adminToken {
			gctx.AbortWithStatusJSON(http.StatusUnauthorized, weberrors.AppError{
				ErrorName:   http.StatusText(http.StatusUnauthorized),
				Description: "invalid admin token",
			})
			return
		}
		gctx.Next()
	}
}
