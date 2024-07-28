package middlewares

import (
	"net/http"
	"scattergories-backend/internal/api/handlers"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func JWTAuthMiddleware(tokenService services.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			handlers.HandleError(c, http.StatusUnauthorized, common.ErrAuthorizationHeaderNotFound.Error())
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if the token is blacklisted
		blacklisted, err := tokenService.IsTokenBlacklisted(tokenString)
		if err != nil {
			handlers.HandleError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		if blacklisted {
			handlers.HandleError(c, http.StatusUnauthorized, common.ErrInvalidToken.Error())
			c.Abort()
			return
		}

		// Check if the token is valid
		claims, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			handlers.HandleError(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Extract and parse the user_id claim
		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			handlers.HandleError(c, http.StatusUnauthorized, "Invalid User ID")
			c.Abort()
			return
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			handlers.HandleError(c, http.StatusUnauthorized, "Invalid UUID format")
			c.Abort()
			return
		}

		// Set the userID and userType in the context
		c.Set("userID", userID)
		c.Set("userType", claims["user_type"])
		c.Next()
	}
}
