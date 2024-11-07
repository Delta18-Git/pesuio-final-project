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
	result := database.DB.Where("ID = ?", request.QuestionID).First(&question)
	if result.Error != nil {
		c.JSON(400, gin.H{
			"error": "question not found",
		})
		return
	} else {
		var testCases []models.TestCase
		err := database.DB.Model(&question).Association("TestCases").Find(&testCases)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "error fetching test cases",
			})
			return
		}
		question.TestCases = testCases
		c.JSON(200, gin.H{
			"success":  true,
			"question": question,
		})
	}

}
