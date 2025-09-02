package questions

import (
	"github.com/delta18-git/taskrunner/database"
	"github.com/delta18-git/taskrunner/models"
	"github.com/gin-gonic/gin"
)

func CreateQuestion(c *gin.Context) {
	var request models.CreateQuestionRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "invalid input",
		})
		return
	}
	question := models.Question{
		Question:  request.Question,
		TestCases: request.TestCases,
		Score:     request.Score,
	}
	result := database.DB.Create(&question)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "unable to create question",
		})
	} else {
		c.JSON(200, gin.H{
			"success":    true,
			"questionID": question.ID,
		})
	}
}
