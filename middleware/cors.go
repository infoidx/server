package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")

		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
			ctx.Header("Access-Control-Allow-Headers", "Authorization")
			ctx.Header("Access-Control-Max-Age", "172800")
			ctx.Header("Access-Control-Allow-Credentials", "true")
			ctx.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		ctx.Next()
	}
}
