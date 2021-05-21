package resource_service

import (
	"easycache/models"
	"easycache/pkg/logger"
	"easycache/pkg/util"
	"easycache/service/config_service"
	"easycache/service/monitor_service"
	"errors"
	"fmt"
	jsonq "github.com/thedevsaddam/gojsonq/v2"
	"strconv"
)

func UpdateAllocatedMemory(info models.ResourceHandle) (err error) {
	resourceInfo, err := GetResourceInfoByHost(info.Ip)
	if err != nil {
		logger.Error("UpdateAllocatedMemory GetResourceInfoByHost error:", err)
		return err
	}
	memory := info.Content
	var newAllocatedMemory int64
	if resourceInfo.AllocatedMemory != "" {
		newAllocatedMemory, err = strconv.ParseInt(resourceInfo.AllocatedMemory, 10, 64)
		logger.Info("newAllocatedMemory before:", newAllocatedMemory)
		if err != nil {
			logger.Error("UpdateAllocatedMemory parseInt error:", err)
			return err
		}
	}
	logger.Info(memory)
	logger.Info(newAllocatedMemory)
	if info.Action == "delete" {
		newAllocatedMemory = newAllocatedMemory - memory/1024/1024
	} else if info.Action == "add" {
		newAllocatedMemory = newAllocatedMemory + memory/1024/1024
	} else {
		return errors.New("memory operation wrong")
	}
	logger.Info("newAllocatedMemory after:", newAllocatedMemory)
	resourceInfo.AllocatedMemory = fmt.Sprintf("%d", newAllocatedMemory)
	err = UpdateResource(resourceInfo)
	if err != nil {
		logger.Error("resource_service UpdateAllocatedMemory err:", err)
		return err
	}
	return nil
}
func InsertResource(resource models.ResourceInfo) (err error) {
	instancesMemory, err := getSumMemory(resource.HostIp)
	if err != nil {
		logger.Error("instance_service InsertResource getSumMemory  err:", err)
		return err
	}
	allocatedUsedMemory, err := getAllocatedUsedMemory(resource.HostIp)
	if err != nil {
		logger.Error("instance_service InsertResource getAllocatedUsedMemory err:", err)
		return err
	}
	resource.AllocatedMemory = fmt.Sprintf("%d", instancesMemory)
	resource.UsedMemory = fmt.Sprintf("%d", allocatedUsedMemory)
	var hostMemory, limitMemory int64
	if resource.TotalHostMemory != "" && resource.LimitedMemory != "" {
		hostMemory, err = strconv.ParseInt(resource.TotalHostMemory, 10, 64)
		if err != nil {
			logger.Error("resource_service InsertResource parseInt hostMemory err:", err)
			return err
		}
		limitMemory, err = strconv.ParseInt(resource.LimitedMemory, 10, 64)
		if err != nil {
			logger.Error("resource_service InsertResource parseInt hostMemory err:", err)
			return err
		}
	}
	if hostMemory >= limitMemory && limitMemory >= instancesMemory || instancesMemory >= allocatedUsedMemory {
		err := models.InsertResourceInfo(resource)
		if err != nil {
			logger.Error("resource_service InsertResource InsertResourceInfo err:", err)
			return err
		}
	} else {
		logger.Error("resource_service InsertResource memory set error:", err)
		logger.Debug("hostMemory", hostMemory, "limitMemory", limitMemory, "instancesMemory", instancesMemory, "allocatedUsedMemory", allocatedUsedMemory)
		return errors.New("resource_service InsertResource memory set error")
	}
	return nil
}

func GetResourceInfoByHost(host string) (resource models.ResourceInfo, err error) {
	resource, err = models.GetResourceInfoByHost(host)
	if err != nil {
		return resource, err
	}
	return resource, nil
}

func GetAllResources() (resources []models.ResourceInfo, err error) {
	resources, err = models.GetAllResources()
	if err != nil {
		return resources, err
	}
	return resources, nil
}

func UpdateResource(resource models.ResourceInfo) error {
	if err := models.UpdateResource(resource); err != nil {
		return err
	}
	return nil
}

func DeleteResource(hostIp string) error {
	if err := models.DeleteResource(hostIp); err != nil {
		return err
	}
	return nil
}

func getSumMemory(hostIp string) (sumMemory int64, err error) {
	instances, err := models.GetActiveInstancesByIp(hostIp)
	if err != nil {
		logger.Error("resource_service getSumMemory err:", err)
		return sumMemory, err
	}
	if len(instances) != 0 {
		for _, instance := range instances {
			config, err := config_service.GetConfigByReplicaId(instance.ID)
			content := config.Content
			maxMemory := jsonq.New().FromString(content).Find("maxmemory.nowValue")
			var value int64
			if maxMemoryStr := util.Interface2String(maxMemory); maxMemoryStr != "" {
				value, err = strconv.ParseInt(maxMemoryStr, 10, 64)
				if err != nil {
					logger.Error("resource_service InsertInstance parseInt err", err)
					return sumMemory, err
				}
			}
			logger.Debug(value)
			sumMemory += value
		}
		// 将B转化为MB
		sumMemory = sumMemory / 1024 / 1024
	}
	return sumMemory, nil
}

func getAllocatedUsedMemory(hostIp string) (usedMemory int64, err error) {
	instances, err := models.GetActiveInstancesByIp(hostIp)
	if err != nil {
		logger.Error("resource_service GetActiveInstancesByIp err:", err)
		return usedMemory, err
	}
	if len(instances) != 0 {
		for _, instance := range instances {
			port, _ := strconv.ParseInt(instance.Port, 10, 64)
			monitorInfo, err := monitor_service.GetMonitorInfo(instance.IP, port, 1)
			if err != nil {
				logger.Error("resource_service getAllocatedUsedMemory err:", err)
				return 0, err
			}
			usedMemoryRaw := jsonq.New().FromString(monitorInfo[0].Info).Find("Memory.used_memory")
			if usedMemoryStr := util.Interface2String(usedMemoryRaw); usedMemoryStr != "" {
				value, err := strconv.ParseInt(usedMemoryStr, 10, 64)
				if err != nil {
					logger.Error("resource_service InsertInstance parseInt err", err)
					return 0, err
				}
				usedMemory += value
			}
		}
		usedMemory = usedMemory / 1024 / 1024
	}
	return usedMemory, nil
}
