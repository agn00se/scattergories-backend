package middleware

import (
	"net/http"
	"scattergories-backend/internal/client/controllers"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			controllers.HandleError(c, http.StatusUnauthorized, common.ErrAuthorizationHeaderNotFound.Error())
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if the token is blacklisted
		blacklisted, err := services.IsTokenBlacklisted(tokenString)
		if err != nil {
			controllers.HandleError(c, http.StatusInternalServerError, err.Error())
			c.Abort()
			return
		}
		if blacklisted {
			controllers.HandleError(c, http.StatusUnauthorized, common.ErrInvalidToken.Error())
			c.Abort()
			return
		}

		// Check if the token is valid
		claims, err := services.ValidateToken(tokenString)
		if err != nil {
			controllers.HandleError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Set user ID and user type in the context
		c.Set("userID", claims["user_id"])
		c.Set("userType", claims["user_type"])
		c.Next()
	}
}
