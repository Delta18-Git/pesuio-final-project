package main

import (
	"github.com/delta18-git/pesuio-final-project/auth"
	"github.com/delta18-git/pesuio-final-project/compiler"
	"github.com/delta18-git/pesuio-final-project/database"
	"github.com/delta18-git/pesuio-final-project/questions"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init("backend.db")
	router := gin.Default()

	router.POST("/auth/signin", auth.Signin)
	router.POST("/auth/signup", auth.Signup)

	router.POST("/run", auth.JwtMiddleware(), compiler.Run)
	router.POST("/testRun", auth.JwtMiddleware(), compiler.RunTest)

	router.POST("/question/create", auth.JwtMiddleware(), questions.CreateQuestion)
	router.POST("/question/fetch", auth.JwtMiddleware(), questions.FetchQuestion)
	router.Run(":1337")
}
