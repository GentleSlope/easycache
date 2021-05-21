package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// NodeConfig ...
type NodeConfig struct {
	InstanceId int64     `form:"instance_id" json:"instance_id"`
	CreatedAt  time.Time `gorm:"autoUpdateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoCreateTime" json:"updated_at"`
	Md5        string    `form:"md5" json:"md5"`
	Content    string    `form:"content" json:"content"`
}

func InsertInstanceConfig(config NodeConfig) (err error) {
	if err := db.Create(&config).Error; err != nil {
		return err
	}
	return nil
}

func GetConfigByReplicaId(id int64) (config NodeConfig, err error) {
	err = db.Where("instance_id = ? ", id).First(&config).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return config, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return config, err
	}

	return config, nil
}

func UpdateConfig(config NodeConfig) (err error) {
	if err = db.Model(&config).Where("instance_id = ?", config.InstanceId).Update("content", config.Content).Error; err != nil {
		return err
	}
	return nil
}
