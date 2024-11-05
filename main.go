package main

import (
	"github.com/delta18-git/pesuio-final-project/auth"
	"github.com/delta18-git/pesuio-final-project/compiler"
	"github.com/delta18-git/pesuio-final-project/questions"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/auth/signin", auth.Signin)
	router.POST("/auth/signup", auth.Signup)

	router.POST("/run", compiler.Run)

	router.POST("/question/create", questions.CreateQuestion)
	router.POST("/question/fetch", questions.FetchQuestion)
	router.Run(":1337")
}
