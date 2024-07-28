package handlers

import (
	"scattergories-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	idStr := c.Param(param)
	return utils.StringToUUID(idStr)
}

func HandleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}
