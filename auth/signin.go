package auth

import (
	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Signin(c *gin.Context) {
	var request models.SignInRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid input",
		})
	}

	ok, err := database.CheckPassword(request.Username, request.Password)
	if ok {
		c.JSON(200, gin.H{"message": "welcome user"})
		return

	} else {
		c.JSON(400, gin.H{"message": "wrong username or password"})
		return
	}

}
