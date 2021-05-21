package models

import (
	"easycache/pkg/logger"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Instance struct {
	ID          int64  `form:"id" gorm:"primary_key" json:"id"`
	RedisUserId int64  `form:"redis_user_id" json:"redis_user_id"`
	Name        string `form:"name" gorm:"name" json:"name"`
	IP          string `form:"ip" json:"ip"`
	Port        string `form:"port" json:"port"`
	IsVirtual   string `form:"is_virtual" json:"is_virtual"`
	IsCluster   string `form:"is_cluster" json:"is_cluster"`
	Password    string `form:"password" json:"password"`
	Status      string `form:"status"  json:"status"`
	Version     string `form:"version"  json:"version"`
}

type ResourceRestrict struct {
	Cpu     int64  `form:"cpu" json:"cpu"`
	Memory  int64  `form:"memory"  json:"memory"`
	Block   int64  `form:"block"  json:"block"`
	ExtInfo string `form:"ext_info"  json:"ext_info"`
}

type InstancePlus struct {
	Instance     Instance         `json:"instance"`
	Uuid         string           `form:"uuid" json:"uuid" `
	Config       NodeConfig       `json:"config"`
	Host         string           `form:"host" json:"host"`
	Port         string           `form:"ex_port" json:"port"`
	UserName     string           `form:"username" json:"username"`
	Password     string           `form:"ex_password" json:"password"`
	IsAutoDeploy string           `form:"is_auto_deploy" json:"is_auto_deploy"`
	Resource     ResourceRestrict `json:"resource"`
}

// InstanceDisplay 用于前端展示
type InstanceDisplay struct {
	Role     string `json:"role"`
	MasterId int64  `json:"master_id"`
	Instance
}

func GetInstancesByUserId(userId int64) (instances []Instance, err error) {
	logger.Debug(userId)
	err = db.Where("redis_user_id = ?", userId).Find(&instances).Error
	logger.Debug(len(instances))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return instances, nil
}

func GetInstancesByIp(hostIp string) (instances []Instance, err error) {
	var instancesTemp []Instance
	err = db.Where("ip= ?", hostIp).Order("id desc").Find(&instancesTemp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	for _, v := range instancesTemp {
		if v.Status != "delete" {
			instances = append(instances, v)
		}
	}
	return instances, nil
}

func GetActiveInstancesByIp(ip string) (instances []Instance, err error) {
	err = db.Where("ip = ? AND status != 'delete' ", ip).Find(&instances).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return instances, nil
}

func GetActiveInstancesBySocket(ip string, port int) (instances []Instance, err error) {
	err = db.Where("ip = ? AND port = ? AND status != 'delete' ", ip, fmt.Sprintf("%d", port)).Find(&instances).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return instances, nil
}

func InsertInstance(instance Instance) (err error) {
	if err := db.Create(&instance).Error; err != nil {
		return err
	}
	return nil
}
func UpdateInstanceWithLastId(instance Instance) (err error) {
	err = db.Where(&instance).Order("id desc").First(&instance).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

func InsertInstanceAndUpdate(instance Instance) (err error) {
	logger.Error(instance)
	if err := db.Create(&instance).Error; err != nil {
		return err
	}
	err = UpdateInstanceWithLastId(instance)
	if err != nil {
		logger.Error("models InsertInstanceAndUpdate err :", err)
		return err
	}
	return nil
}

func GetInstanceById(id int64) (instance Instance, err error) {
	err = db.Where("id = ?", id).First(&instance).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return instance, err
	}

	return instance, nil
}

func GetInstances() (instances []Instance, err error) {
	err = db.Where("status != 'delete'").Find(&instances).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return instances, nil
}

func UpdateTaskStatus(instance Instance) (err error) {
	logger.Info("UpdateTaskStatus", instance)
	err = db.Model(&instance).Where("id = ?", instance.ID).Update("status", instance.Status).Error
	if err != nil {
		return err
	}
	return nil
}
