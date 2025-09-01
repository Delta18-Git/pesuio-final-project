package compiler

import (
	"net/http"
	"os"
	"strings"

	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Run(c *gin.Context) {
	var request models.RunRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request payload",
		})
		return
	}
	tempFile, err := os.CreateTemp("", "code-*."+request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create temporary file to store code",
		})
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(request.Code))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to write code to temp file",
		})
		return
	}
	output, errors := runDocker(tempFile, request.Language, request.Input)
	if output == "" || errors != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"output":  output,
			"errors":  errors,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"output":  output,
			"errors":  errors,
		})
	}

}

func RunTest(c *gin.Context) {
	var request models.TestRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request payload",
		})
		return
	}
	var question models.Question
	result := database.DB.Preload("TestCases").Where("ID = ?", request.QuestionID).First(&question)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "question not found",
		})
		return
	}
	tempFile, err := os.CreateTemp("", "code-*."+request.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to create temporary file to store code",
		})
		return
	}
	defer os.Remove(tempFile.Name())
	_, err = tempFile.Write([]byte(request.Code))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "failed to write code to temp file",
		})
		return
	}
	for _, testCase := range question.TestCases {
		input := testCase.Input
		output, errors := runDocker(tempFile, request.Language, input)
		if testCase.ExpectedOutput != strings.Trim(output, "\n ") || errors != "" {
			c.JSON(http.StatusExpectationFailed, gin.H{
				"success":        false,
				"output":         output,
				"expectedOutput": testCase.ExpectedOutput,
				"errors":         errors,
			})
			return
		}
	}
	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
		"message": "all test cases passed",
	})
}
