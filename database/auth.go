package database

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/delta18-git/pesuio-final-project/models"
)

func CreateUser(username, password string) error {
	// creates a new user in the database, returns error if any
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := models.User{
		Username: username,
		Password: string(hashedPw),
	}
	return DB.Create(&newUser).Error
}

func CheckPassword(username, password string) (success bool, err error) {
	// checks if the password is correct for the given username
	var checkDB models.User
	result := DB.Where("Username = ?", username).First(&checkDB)
	if result.Error != nil {
		return false, result.Error
	}
	err = bcrypt.CompareHashAndPassword([]byte(checkDB.Password), []byte(password))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func CheckUser(username string) (success bool, err error) {
	// checks if the given username already exists
	var checkDB models.User
	result := DB.Where("Username = ?", username).First(&checkDB)
	if result.Error != nil {
		return false, result.Error
	} else {
		return true, nil
	}
}
