package auth

import (
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")
		splitToken := strings.Split(bearer, " ") // remove the "Bearer <token>" to just <token>
		tokenString := ""
		if len(splitToken) == 2 {
			tokenString = splitToken[1]
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			key, found := os.LookupEnv("JWT_KEY")
			if !found {
				return []byte(""), errors.New("missing jwt signing key")
			}
			return []byte(key), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}))
		if err != nil {
			c.JSON(401, gin.H{"success": false, "message": "unable to authenticate"})
			c.AbortWithError(401, errors.New("unable to parse JWT token"))
			return
		}
		_, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(401, gin.H{"success": false, "message": "authentication required, not authorized"})
			return
		}
		c.Next()
	}
}
