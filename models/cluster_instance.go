package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type ClusterInstance struct {
	ClusterId  int64     `gorm:"primary_key" json:"cluster_id"`
	InstanceId int64     `gorm:"primary_key" json:"instance_id"`
	Slots      string    `json:"slots"`
	Role       string    `json:"role"`
	MasterId   int64     `json:"master_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func GetClusterInstanceById(clusterId int64, instanceId int64) (isExist bool, err error) {
	var clusterInstances []ClusterInstance
	err = db.Where("cluster_id = ? AND instance_id = ?", clusterId, instanceId).Find(&clusterInstances).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return isExist, err
	}

	if len(clusterInstances) > 0 {
		isExist = true
	}
	return isExist, nil
}

func GetClusterInstancesById(clusterId int64) (clusterInstances []ClusterInstance, err error) {
	err = db.Where("cluster_id = ?", clusterId).Find(&clusterInstances).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return clusterInstances, err
	}

	return clusterInstances, nil
}

func GetClusterInstancesByInstanceId(instanceId int64) (clusterInstance ClusterInstance, err error) {
	err = db.Where("instance_id = ?", instanceId).First(&clusterInstance).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return clusterInstance, err
	}

	return clusterInstance, nil
}

func InsertClusterInstance(clusterInstance ClusterInstance) (err error) {
	if err := db.Create(&clusterInstance).Error; err != nil {
		return err
	}
	return nil
}

func DeleteClusterInstance(clusterInstance ClusterInstance) (err error) {
	if err := db.Where("instance_id = ?", clusterInstance.InstanceId).Delete(&clusterInstance).Error; err != nil {
		return err
	}
	return nil
}
