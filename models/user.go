package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(20);not null;"`
	Password string `gorm:"type:varchar(50);not null;"`
	PhoneNum string `gorm:"type:char(11);not null;"`
}

// AddUser add a user
func AddUser(username, password, phoneNum string) (uint, error) {
	user := User{
		Username: username,
		Password: password,
		PhoneNum: phoneNum,
	}
	if err := db.Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// ExistUserByName determine whether a user exists based on username
func ExistUserByName(username string) (bool, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	return user.ID > 0, nil
}

// ExistUserByPhone determine whether a user exists based on phoneNum and meanwhile get ID
func ExistUserByPhone(phoneNum string) (bool, uint, error) {
	var user User
	err := db.Where("phone_num = ?", phoneNum).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, 0, err
	}

	return user.ID > 0, user.ID, nil
}

// CheckUser get a user's ID if providing username and password matches
func CheckUser(username, password string) (uint, error) {
	var user User
	err := db.Where(&User{Username: username, Password: password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	return user.ID, nil
}

// DeleteUser delete a user by ID (soft delete)
func DeleteUser(id uint) error {
	var user User
	if err := db.Delete(&user, id).Error; err != nil {
		return err
	}

	return nil
}

// GetUser get user infomation by ID
func GetUser(userID uint) (User, error) {
	var user User
	err := db.Where("id = ?", userID).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return User{}, err
	}

	return user, nil
}
