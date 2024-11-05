package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

type Question struct {
	gorm.Model
	Question  string     `json:"question"`
	TestCases []TestCase `json:"testCases"`
	Score     int        `json:"score"`
}

type TestCase struct {
	gorm.Model
	QuestionID     uint
	Input          string `json:"input"`
	ExpectedOutput string `json:"expectedOutput"`
}
