package controllers

import (
	"net/http"
	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/client/controllers/responses"
	"strings"

	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := services.Login(req.Email, req.Password)
	if err != nil {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}
	response := responses.ToLoginResponse(accessToken, refreshToken)
	c.JSON(http.StatusOK, response)
}

func Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if err := services.InvalidateToken(tokenString); err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func ExchangeToken(c *gin.Context) {
	var req requests.ExchangeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	newAccessToken, newRefreshToken, err := services.RefreshTokens(req.RefreshToken)
	if err != nil {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	response := responses.ToExchangeTokenResponse(newAccessToken, newRefreshToken)
	c.JSON(http.StatusOK, response)
}
