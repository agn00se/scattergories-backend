package controllers

import (
	"net/http"

	"scattergories-backend/internal/client/controllers/requests"
	"scattergories-backend/internal/client/controllers/responses"
	"scattergories-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	users, err := services.GetAllUsers()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve users")
	}

	var response []responses.UserResponse
	// _: The blank identifier _ is used to ignore the index of the slice or array.
	for _, user := range users {
		response = append(response, responses.ToUserResponse(user))
	}

	c.JSON(http.StatusOK, response)
}

func GetUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := services.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "User not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to get user")
		}
		return
	}

	response := responses.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

func CreateUser(c *gin.Context) {
	user, err := services.CreateGuestUser()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create user")
		return
	}

	response := responses.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

func UpdateUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var request requests.UserRequest
	// short-lived variables like err can be made local within the if statement
	if err := c.ShouldBindJSON(&request); err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := services.UpdateUserByID(id, request.Name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "User not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to update user")
		}
		return
	}

	response := responses.ToUserResponse(user)
	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = services.DeleteUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			handleError(c, http.StatusNotFound, "User not found")
		} else {
			handleError(c, http.StatusInternalServerError, "Failed to delete user")
		}
		return
	}

	// gin.H{} is a shortcut provided by Gin for map[string]interface{}.
	// It is used to simplify the creation of JSON responses.
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
