package compiler

import (
	"net/http"
	"os"

	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Run(c *gin.Context) {
	var request models.RunRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	tempFile, err := os.CreateTemp("", "code-*."+request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file to store code"})
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(request.Code))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write code to temp file"})
		return
	}
	output, errors := runDocker(tempFile, request.Language, request.Input)
	c.JSON(http.StatusOK, gin.H{
		"output": output + errors,
	})
}
