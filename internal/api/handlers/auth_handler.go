package handlers

import (
	"net/http"
	"scattergories-backend/internal/api/handlers/requests"
	"scattergories-backend/internal/api/handlers/responses"
	"strings"

	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	Login(c *gin.Context)
	Logout(c *gin.Context)
	ExchangeToken(c *gin.Context)
}

type AuthHandlerImpl struct {
	authService  services.AuthService
	tokenService services.TokenService
}

func NewAuthHandler(authService services.AuthService, tokenService services.TokenService) AuthHandler {
	return &AuthHandlerImpl{authService: authService, tokenService: tokenService}
}

func (h *AuthHandlerImpl) Login(c *gin.Context) {
	var req requests.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}
	response := responses.ToLoginResponse(accessToken, refreshToken)
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandlerImpl) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if err := h.tokenService.InvalidateToken(tokenString); err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandlerImpl) ExchangeToken(c *gin.Context) {
	var req requests.ExchangeTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	newAccessToken, newRefreshToken, err := h.tokenService.RefreshTokens(req.RefreshToken)
	if err != nil {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	response := responses.ToExchangeTokenResponse(newAccessToken, newRefreshToken)
	c.JSON(http.StatusOK, response)
}
