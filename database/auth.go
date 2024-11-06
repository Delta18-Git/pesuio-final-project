package database

import "github.com/delta18-git/pesuio-final-project/models"

func CreateUser(username, password string) error {
	// creates a new user in the database, returns error if any
	newUser := models.User{
		Username: username,
		Password: password,
	}
	return DB.Create(&newUser).Error
}

func CheckPassword(username, password string) (success bool, err error) {
	// checks if the password is correct for the given username
	var checkDB models.User
	result := DB.Where("Username = ? AND Password = ?", username, password).First(&checkDB)
	if result.Error != nil {
		return false, result.Error
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