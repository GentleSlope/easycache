package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type RedisUser struct {
	ID        int64     `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Uuid      string    `json:"uuid"`
	Details   string    `json:"details"`
	CreatedAt time.Time `gorm:"autoUpdateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoCreateTime" json:"updated_at"`
}

func GetRedisUserById(redisUserId int64) (user RedisUser, err error) {
	err = db.Where("id = ? ", redisUserId).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, nil
}

func GetRedisUsers() (users []RedisUser, err error) {
	err = db.Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return users, err
	}
	return users, nil
}

func FindRedisUserByName(name string) (user RedisUser, err error) {
	err = db.Where("name = ? ", name).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}

	return user, nil
}

func InsertRedisUser(user *RedisUser) (err error) {
	if err = db.Create(&user).Error; err != nil {
		return err
	}
	*user, err = FindRedisUserByName(user.Name)
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserById(id int64) (err error) {
	if err := db.Where("id = ?", id).Delete(&RedisUser{}).Error; err != nil {
		return err
	}
	return nil
}
