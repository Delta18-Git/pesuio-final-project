package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Signin(c *gin.Context) {
	var request models.SignInRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid input",
			"success": false,
		})
		return
	}

	key, found := os.LookupEnv("JWT_KEY")
	if !found {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error, login not available",
			"success": false,
		})
		return
	}
	ok, err := database.CheckPassword(request.Username, request.Password)
	if ok {
		claims := jwt.MapClaims{}
		claims["authorized"] = true
		claims["exp"] = time.Now().Add(time.Hour * 12).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		signedToken, signErr := token.SignedString([]byte(key))
		if signErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "internal server error, login not available",
			})
		} else {
			// not using cookie as there is no reference frontend
			c.JSON(http.StatusOK, gin.H{
				"message":      "welcome user",
				"success":      true,
				"access_token": signedToken,
			})
		}
		return

	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong username or password",
			"success": false,
		})
		return
	}

}
