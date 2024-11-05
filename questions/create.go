package questions

import (
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func CreateQuestion(c *gin.Context) {
	var request models.CreateQuestionRequest
	c.BindJSON(&request)

	c.JSON(200, gin.H{
		"success":    true,
		"questionID": 1,
	})

}
