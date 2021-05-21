package models

import "github.com/jinzhu/gorm"

type ResourceInfo struct {
	HostIp            string `gorm:"primary_key" form:"host_ip" json:"host_ip"`
	HostGroupIdentity string `form:"host_group_identity" json:"host_group_identity"`
	TotalHostMemory   string `form:"total_host_memory" json:"total_host_memory"`
	UsedMemory        string `form:"used_memory" json:"used_memory"`
	AllocatedMemory   string `form:"allocated_memory" json:"allocated_memory"`
	LimitedMemory     string `form:"limited_memory" json:"limited_memory"`
}

type ResourceHandle struct {
	Ip      string `json:"ip"`
	Content int64  `json:"content"`
	Action  string `json:"action"`
}

func GetResourceInfoByHost(host string) (resource ResourceInfo, err error) {
	err = db.Where("host_ip = ? ", host).First(&resource).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return resource, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return resource, err
	}

	return resource, nil
}

func GetAllResources() (resources []ResourceInfo, err error) {
	err = db.Find(&resources).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return resources, err
	}

	return resources, nil
}

func InsertResourceInfo(resource ResourceInfo) (err error) {
	if err := db.Create(&resource).Error; err != nil {
		return err
	}
	return nil
}

func UpdateResource(resource ResourceInfo) (err error) {
	if err = db.Save(&resource).Error; err != nil {
		return err
	}
	return nil
}
func DeleteResource(hostIp string) (err error) {
	if err := db.Where("host_ip = ?", hostIp).Delete(&ResourceInfo{}).Error; err != nil {
		return err
	}
	return nil
}
