package auth

import (
	"net/http"

	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var request models.SignUpRequest
	err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid input",
		})
		return
	}
	ok, _ := database.CheckUser(request.Username)
	if !ok {
		database.CreateUser(request.Username, request.Password)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "user already exists",
		})
	}
}
