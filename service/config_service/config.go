package config_service

import (
	"easycache/models"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/task_service"
)

func GetConfigByReplicaId(id int64) (config models.NodeConfig, err error) {
	config, err = models.GetConfigByReplicaId(id)
	if err != nil {
		return config, err
	}
	return config, nil
}

func UpdateConfig(config models.NodeConfig) (err error) {
	err = models.UpdateConfig(config)
	if err != nil {
		logger.Error("config_service UpdateConfig err:", err)
		return err
	}
	task := models.Task{ReplicaId: config.InstanceId, TaskDetails: config.Content, TaskStatus: define.TaskWaiting, TaskType: define.TaskChangeConfig}
	err = task_service.ADD(task)
	if err != nil {
		logger.Error("config_service add task err:", err)
		return err
	}
	return nil
}
