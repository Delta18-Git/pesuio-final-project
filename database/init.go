package database

import (
	"github.com/delta18-git/pesuio-final-project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbErr error

func Init(databaseFileName string) {
	DB, dbErr = gorm.Open(sqlite.Open("backend.db"), &gorm.Config{})
	if dbErr != nil {
		panic("Error connecting to database")
	}
	// implement
	DB.AutoMigrate(&models.User{}, &models.Question{})

}
