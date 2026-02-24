package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	gosiftjwt "github.com/JoseIgnacioGC/gosift-backend/internal/platform/jwt"
)

const UserIDKey = "userID"

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		token, found := strings.CutPrefix(header, "Bearer ")
		if !found {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format, expected 'Bearer <token>'"})
			return
		}

		claims, err := gosiftjwt.ValidateToken(token, jwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}
