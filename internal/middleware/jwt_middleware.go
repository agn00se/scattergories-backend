package middleware

import (
	"net/http"
	"scattergories-backend/internal/client/controllers"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"
	"strconv"
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

		// Check if the user ID from the request matches the user ID from the token
		userIDFromToken := uint(claims["user_id"].(float64))
		userIDFromRequest := c.DefaultPostForm("user_id", c.DefaultPostForm("host_id", ""))

		if userIDFromRequest != "" {
			userIDConverted, err := strconv.Atoi(userIDFromRequest)
			if err != nil || uint(userIDConverted) != userIDFromToken {
				controllers.HandleError(c, http.StatusUnauthorized, common.ErrInvalidToken.Error())
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
