package auth

import (
	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var request models.SignUpRequest
	err := c.BindJSON(&request)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid input",
		})
		return
	}
	ok, _ := database.CheckUser(request.Username)
	if !ok {
		database.CreateUser(request.Username, request.Password)
		c.JSON(200, gin.H{
			"success": true,
		})
	} else {
		c.JSON(400, gin.H{
			"error": "user already exists",
		})
	}
}
