package models

type SignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RunRequest struct {
	Language string `json:"language" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Input    string `json:"input"`
}

type TestRequest struct {
	Language   string `json:"language" binding:"required"`
	Code       string `json:"code" binding:"required"`
	QuestionID uint   `json:"questionID" binding:"required"`
}

type CreateQuestionRequest struct {
	Question  string     `json:"question" binding:"required"`
	TestCases []TestCase `json:"testCases" binding:"required"`
	Score     int        `json:"score" binding:"required"`
}

type FetchQuestionRequest struct {
	QuestionID uint `json:"questionID" binding:"required"`
}
