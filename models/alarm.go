package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type AlarmInfo struct {
	RedisUserId int64     `json:"redis_user_id"`
	InstanceID  int64     `json:"instance_id"`
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Msg         string    `json:"msg"`
	CreatedAt   time.Time `gorm:"autoUpdateTime" json:"created_at"`
}

func InsertAlarmInfo(info AlarmInfo) error {
	if err := db.Create(&info).Error; err != nil {
		return err
	}
	return nil
}

func GetAllAlarmInfos() (infos []AlarmInfo, err error) {
	err = db.Order("id desc").Find(&infos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return infos, err
	}
	return infos, nil
}

func DeleteAlarmInfos() (err error) {
	if err := db.Where("id != 0").Delete(&AlarmInfo{}).Error; err != nil {
		return err
	}
	return nil
}
