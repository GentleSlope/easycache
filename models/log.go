package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type LogInfo struct {
	UUID         string    `json:"uuid"`
	RedisUserId  string    `json:"redis_user_id"`
	InstanceIP   string    `json:"instance_ip"`
	InstancePort string    `json:"instance_port"`
	LogContent   string    `json:"log_content"`
	LogLevel     string    `json:"log_level"`
	CreatedAt    time.Time `gorm:"autoUpdateTime" json:"created_at"`
}

func InsertLogs(log LogInfo) error {
	if err := db.Create(&log).Error; err != nil {
		return err
	}
	return nil
}

func GetAllLogs(ip string, port string, limit int) (logs []LogInfo, err error) {
	err = db.Limit(limit).Where(" instance_ip = ? AND instance_port = ?", ip, port).Order("id desc").Find(&logs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return logs, err
	}
	return logs, nil
}

func DeleteLogs(ip string, port string) (err error) {
	if err := db.Where("instance_ip = ? AND instance_port = ?", ip, port).Delete(&LogInfo{}).Error; err != nil {
		return err
	}
	return nil
}
