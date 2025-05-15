package handlers

import (
	"hospital-management/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader || tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization format"})
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if role, ok := claims["role"].(string); ok && role == requiredRole {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
	}
}
