package models

import (
	"easycache/pkg/logger"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Task struct {
	ID              int64     `gorm:"primary_key" json:"id"`
	ReplicaId       int64     `gorm:"index" json:"replica_id"`
	TaskType        string    `json:"task_type"`
	TaskDetails     string    `json:"task_details"`
	TaskStatus      string    `json:"task_status"`
	ExecutionResult string    `json:"execution_result"`
	Instance        Instance  `json:"instance"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// GetTask Get a single task based on ReplicaId
func GetTasksByReplicaId(replicaId int64) (tasks []Task, err error) {
	err = db.Where("replica_id = ?", replicaId).Find(&tasks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tasks, nil
}

func GetTaskByReplicaId(id int64) (task Task, err error) {
	err = db.Where("replica_id = ? ", id).First(&task).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return task, err
	}
	return task, nil
}

func GetTaskById(id int64) (task Task, err error) {
	err = db.Where("id = ? ", id).First(&task).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return task, err
	}
	return task, nil
}

func GetTasksBetween(instanceIds []int64, start string, end string) (tasks []Task, err error) {
	logger.Debug(instanceIds)
	var temp = make([]string, len(instanceIds))
	for k, v := range instanceIds {
		temp[k] = strconv.FormatInt(v, 10)
	}
	logger.Debug(temp)
	err = db.Where("replica_id IN ( ? ) AND created_at > ? AND updated_at < ?", temp, start, end).Find(&tasks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tasks, nil
}

func GetWaitingDeployTasks() (tasks []Task, err error) {
	err = db.Where("task_status = 'waiting' AND task_type = 'deploy' ").Order("id desc").Find(&tasks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tasks, nil
}

func GetRunningNotDeployTasks(replicaId int64) (tasks []Task, err error) {
	err = db.Where("task_status = 'running' AND task_type != 'deploy' AND replica_id = ?", replicaId).Order("id desc").Find(&tasks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tasks, nil
}

func GetWaitingNotDeployTask(replicaId int64) (task Task, err error) {
	err = db.Where("task_status = 'waiting' AND task_type != 'deploy' AND replica_id = ?", replicaId).Order("id asc").First(&task).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return task, err
	}
	return task, nil
}

func UpdateTask(task Task) (err error) {
	if err = db.Save(&task).Error; err != nil {
		return err
	}
	return nil
}

// AddTask ...
func AddTask(data Task) error {
	if err := db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTaskById(id int64) error {
	if err := db.Where("id = ?", id).Delete(&Task{}).Error; err != nil {
		return err
	}
	return nil
}
