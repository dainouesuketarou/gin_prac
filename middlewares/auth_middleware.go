package middlewares

import (
	"first_gin_app/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.IAuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader(("Authorization"))
		if header == "" {
			// 401 Unauthorizedステータスで終了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			// 401 Unauthorizedステータスで終了
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		user, err := authService.GetUserFromToken(tokenString)
		if err != nil {

			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// ginフレームワークのコンテキストのuserをセット
		ctx.Set("user", user)

		// 次のハンドラに進む
		ctx.Next()
	}
}
