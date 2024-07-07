package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func getIDParam(c *gin.Context, param string) (uint, error) {
	// strconv.Atoi - ASCII to integer
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func handleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
