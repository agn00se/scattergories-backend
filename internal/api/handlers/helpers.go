package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIDParam(c *gin.Context, param string) (uint, error) {
	// strconv.Atoi - ASCII to integer
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func HandleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
