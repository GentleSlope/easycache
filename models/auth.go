package models

import (
	"github.com/jinzhu/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Limit    int    `json:"limit"`
	Email    string `json:"email"`
}

// CheckAuth checks if authentication information exists
func CheckAuth(username string, password string) (bool, Auth, error) {
	var auth Auth
	err := db.Where(Auth{Username: username, Password: password}).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, auth, err
	}

	if auth.ID > 0 && auth.Limit > 0 {
		return true, auth, nil
	}

	return false, auth, nil
}

// ...
func GetPassword(username string) (pwd string, err error) {
	var auth Auth
	err = db.Where("username = ? ", username).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth.Password, err
	}
	return auth.Password, nil
}

// ...
func GetUser(username string) (user Auth, err error) {
	var auth Auth
	err = db.Where("username = ? ", username).First(&auth).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return auth, err
	}
	auth.Password = ""
	return auth, nil
}

// ...
func Register(username string, password string) (err error) {
	var auth Auth
	auth.Email = "1419418290@qq.com"
	auth.Username = username
	auth.Password = password
	if err := db.Create(&auth).Error; err != nil {
		return err
	}
	return nil
}
