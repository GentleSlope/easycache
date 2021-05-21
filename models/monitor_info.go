package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type MonitorInfo struct {
	Id          int64     `json:"id"`
	RedisUserId int64     `json:"redis_user_id"`
	ReplicaIp   string    `json:"replica_ip"`
	ReplicaPort int64     `json:"replica_port"`
	Info        string    `json:"info"`
	CreatedAt   time.Time `json:"created_at"`
}

func InsertInfo(info MonitorInfo) (err error) {
	if err := db.Create(&info).Error; err != nil {
		return err
	}
	return nil
}

func GetMonitorInfo(ip string, port int64, limit int) (infos []MonitorInfo, err error) {
	err = db.Limit(limit).Where(" replica_ip = ? AND replica_port = ?", ip, port).Order("id desc").Find(&infos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return infos, err
	}
	return infos, nil
}
