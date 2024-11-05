package auth

import (
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
	}

	// implement
	c.JSON(200, gin.H{
		"success": true,
	})
}
