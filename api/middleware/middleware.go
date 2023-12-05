package middleware

import (
	"net/http"
	"public-surf/pkg/jwttoken"
	"public-surf/pkg/response"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Claims struct {
	Email string `json:"user_email"`
	jwt.StandardClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var jwtSecret = []byte(viper.GetString("Jwt.Secret"))

		err := jwttoken.TokenValid(c.Request)
		if err != nil {
			response.ResponseError(c, err.Error(), http.StatusUnauthorized)
			c.Abort()
			return
		}

		// decode
		// claims, err := jwttoken.ExtractTokenMetadata(c.Request)
		parts := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		tokenString := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Attach the user information to the context
		c.Set("user_email", claims.Email)

		// Continue to the next handler
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
