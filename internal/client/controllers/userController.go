package controllers

import (
	"net/http"

	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/common"
	"scattergories-backend/internal/models"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
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

func GetUser(c *gin.Context) {
	id, err := GetIDParam(c, "id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := services.GetUserByID(id)
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

func CreateUser(c *gin.Context) {
	var request requests.UserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	var err error
	var user *models.User

	if request.Type == string(models.UserTypeGuest) {
		user, err = services.CreateGuestUser()
	} else if request.Type == string(models.UserTypeRegistered) {
		user, err = services.CreateRegisteredUser(*request.Name, *request.Email, *request.Password)
	}

	if err != nil {
		if err == common.ErrEmailAlreadyUsed {
			HandleError(c, http.StatusConflict, err.Error())
		}
		HandleError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := responses.ToUserResponse(user)
	c.JSON(http.StatusCreated, response)
}

func DeleteUser(c *gin.Context) {
	id, err := GetIDParam(c, "id")
	if err != nil {
		HandleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = services.DeleteUserByID(id)
	if err != nil {
		if err == common.ErrUserNotFound {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	// gin.H{} is a shortcut provided by Gin for map[string]interface{}.
	// It is used to simplify the creation of JSON responses.
	c.JSON(http.StatusNoContent, gin.H{"message": "User deleted"})
}
