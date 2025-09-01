package questions

import (
	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func FetchQuestion(c *gin.Context) {
	var request models.FetchQuestionRequest
	c.BindJSON(&request)

	var question models.Question
	result := database.DB.Preload("TestCases").Where("ID = ?", request.QuestionID).First(&question)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "question not found",
		})
		return
	} else {
		c.JSON(200, gin.H{
			"success":  true,
			"question": question,
		})
	}

}
