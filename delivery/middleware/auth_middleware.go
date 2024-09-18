package middleware

import (
	"fmt"
	"polen/utils/security"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authHeader authHeader
		if err := c.ShouldBindHeader(&authHeader); err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": fmt.Sprintf("unauthorized %v", err.Error()),
			})
			return
		}

		// Substring
		token := strings.Replace(authHeader.AuthorizationHeader, "Bearer ", "", 1)
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "unauthorized",
			})
			return
		}

		// Verifikasi Token
		claims, err := security.VerifyJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": fmt.Sprintf("unauthorized %v", err.Error()),
			})
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}
