package utils

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIntParam(c *gin.Context, paramName string) (int, bool) {

	paramStr := c.Param(paramName)

	// Convert the string value to an integer
	paramInt, err := strconv.Atoi(paramStr)
	if err != nil {
		// Return an error response if the conversion fails
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid " + paramName + " format. It should be an integer.",
		})
		return 0, false // Return 0 and false indicating failure
	}

	return paramInt, true
}
