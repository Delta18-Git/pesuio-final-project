package auth

import (
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func Signin(c *gin.Context) {
	var request models.SignInRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid input",
		})
	}
	
	result := DB.Where("Username = ? AND Password = ?", request.Username, request.Password).First(&request)
	if result.Error == nil{
		c.JSON(200, gin.H{"message":"welcome user"} )
		return


	}
	if result.Error != nil{
		c.JSON(400, gin.H{"message":"wrong username or password"} )
		return


	}

	

}
