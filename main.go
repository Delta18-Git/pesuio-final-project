package main

import (
	"github.com/delta18-git/taskrunner/auth"
	"github.com/delta18-git/taskrunner/compiler"
	"github.com/delta18-git/taskrunner/database"
	"github.com/delta18-git/taskrunner/questions"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Init("backend.db")
	router := gin.Default()

	router.POST("/auth/signin", auth.Signin)
	router.POST("/auth/signup", auth.Signup)

	router.POST("/run/code", auth.JwtMiddleware(), compiler.Run)
	router.POST("/run/tests", auth.JwtMiddleware(), compiler.RunTest)

	router.POST("/question/create", auth.JwtMiddleware(), questions.CreateQuestion)
	router.POST("/question/fetch", auth.JwtMiddleware(), questions.FetchQuestion)
	router.Run(":1337")
}
