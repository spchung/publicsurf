package middleware

import (
	"net/http"
	"public-surf/pkg/jwttoken"
	"public-surf/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := jwttoken.TokenValid(c.Request)
		if err != nil {
			response.ResponseError(c, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}

func IsSignedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := jwttoken.ExtractTokenMetadata(c.Request)
		if err != nil {
			response.ResponseError(c, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
