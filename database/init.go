package database

import (
	"github.com/delta18-git/taskrunner/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbErr error

func Init(databaseFileName string) {
	DB, dbErr = gorm.Open(sqlite.Open(databaseFileName), &gorm.Config{})
	if dbErr != nil {
		panic("Error connecting to database")
	}
	DB.AutoMigrate(&models.User{}, &models.Question{}, &models.TestCase{})

}
