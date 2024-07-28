package handlers

import (
	"net/http"
	"strings"

	"scattergories-backend/internal/api/handlers/requests"
	"scattergories-backend/internal/api/handlers/responses"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	GetAllUsers(c *gin.Context)
	GetUser(c *gin.Context)
	CreateAccount(c *gin.Context)
	CreateGuestAccount(c *gin.Context)
	DeleteAccount(c *gin.Context)
}

type UserHandlerImpl struct {
	userService             services.UserService
	userRegistrationService services.UserRegistrationService
	tokenService            services.TokenService
	permissionService       services.PermissionService
}

func NewUserHandler(
	userService services.UserService,
	userRegistrationService services.UserRegistrationService,
	tokenService services.TokenService,
	permissionService services.PermissionService) UserHandler {
	return &UserHandlerImpl{
		userService:             userService,
		userRegistrationService: userRegistrationService,
		tokenService:            tokenService,
		permissionService:       permissionService,
	}
}

func (h *UserHandlerImpl) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to retrieve users")
	}

	var response []*responses.UserResponse
	// _: The blank identifier _ is used to ignore the index of the slice or array.
	for _, user := range users {
		response = append(response, responses.ToUserResponse(user))
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandlerImpl) GetUser(c *gin.Context) {
	id, err := GetIDParam(c, "id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if err == common.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to get user")
		}
		return
	}

	response := responses.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandlerImpl) CreateAccount(c *gin.Context) {
	var request requests.UserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userRegistrationService.CreateRegisteredUser(request.Name, request.Email, request.Password)
	if err != nil {
		if err == common.ErrEmailAlreadyUsed {
			HandleError(c, http.StatusConflict, err.Error())
			return
		}
		HandleError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := responses.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandlerImpl) CreateGuestAccount(c *gin.Context) {
	user, err := h.userService.CreateGuestUser()
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	// Generate access token for guest users during account creation
	accessToken, err := h.tokenService.GenerateJWT(user.ID, user.Type)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
	}

	response := responses.ToGuestUserResponse(user, accessToken)
	c.JSON(http.StatusCreated, response)
}

func (h *UserHandlerImpl) DeleteAccount(c *gin.Context) {
	id, err := GetIDParam(c, "id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userID, _ := c.Get("userID")
	permitted, err := h.permissionService.HasPermission(userID.(uint), services.UserWritePermission, id)
	if err != nil || !permitted {
		HandleError(c, http.StatusUnauthorized, err.Error())
		return
	}

	err = h.userService.DeleteUserByID(id)
	if err != nil {
		if err == common.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	tokenString := c.GetHeader("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	if err := h.tokenService.InvalidateToken(tokenString); err != nil {
		HandleError(c, http.StatusInternalServerError, "Failed to invalidate token")
		return
	}

	// gin.H{} is a shortcut provided by Gin for map[string]interface{}.
	// It is used to simplify the creation of JSON responses.
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted"})
}
