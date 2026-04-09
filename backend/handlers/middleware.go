package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func jwtAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		split := strings.Split(bearerToken, " ")
		if len(split) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := parseToken(split[1], secret)
		if err != nil {
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Signature"})
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is either expired or not yet valid"})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			}
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, claims.UID)
		c.Next()
	}
}
