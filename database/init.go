package database

import ("gorm.io/gorm"
		"gorm.io/driver/sqlite"
	)

var DB *gorm.DB

func Init(databaseFileName string) {
	DB, err = gorm.Open(sqlite.Open("backend.db"), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}
	// implement
	// populate DB variable	

}
