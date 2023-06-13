package middleware

import (
	"net/http"

	"github.com/anodazDev/go_dev/service"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get the JWT token from the cookie
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Validate the JWT token
		token, err := service.JWTAuthService().ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		// Pass the validated token to the next handler
		c.Set("token", token)
		c.Next()

	}
}
