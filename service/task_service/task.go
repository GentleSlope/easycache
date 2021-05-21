package task_service

import (
	"easycache/models"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"time"
)

func Get(replicaId int64) (tasks []models.Task, err error) {
	tasks, err = models.GetTasksByReplicaId(replicaId)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetNodeChangeTask(replicaId int64) (taskWaiting models.Task, err error) {
	// 查找非部署类型的waiting任务
	waitingTasks, err := models.GetRunningNotDeployTasks(replicaId)
	if err != nil {
		logger.Error("task_service GetNodeChangeTask GetRunningNotDeployTasks err", err)
		return taskWaiting, err
	}
	// 如果当前队列中存在running的任务,有则判断任务是否超时,如果超时则将任务状态设置为failed，若没超时则返回null
	if len(waitingTasks) > 0 {
		for _, task := range waitingTasks {
			updatedAt := task.UpdatedAt
			timeNow := time.Now()
			duringTime := timeNow.Sub(updatedAt)
			// 设置一小时的超时间
			if duringTime > time.Duration(3600*time.Second) {
				task.TaskStatus = define.TaskFailed
				err = UpdateTask(task)
				if err != nil {
					logger.Error("task_service GetNodeChangeTask UpdateTask err:", err)
					return taskWaiting, err
				}
			} else {
				// 有运行的任务，就先不拉取了
				return taskWaiting, nil
			}
		}
	}
	taskWaiting, err = models.GetWaitingNotDeployTask(replicaId)
	if err != nil {
		logger.Error("task_service GetNodeChangeTask GetWaitingNotDeployTasks err", err)
		return taskWaiting, err
	}
	return taskWaiting, nil
}

func GetTaskById(id int64) (task models.Task, err error) {
	task, err = models.GetTaskById(id)
	if err != nil {
		return task, err
	}
	return task, nil
}

func GetTaskByReplicaId(id int64) (task models.Task, err error) {
	task, err = models.GetTaskByReplicaId(id)
	if err != nil {
		return task, err
	}
	return task, nil
}

func GetBetweenData(redisUserId int64, start string, end string) (tasks []models.Task, err error) {
	instances, err := models.GetInstancesByUserId(redisUserId)
	if err != nil {
		return nil, err
	}
	var instanceIds []int64
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.ID)
	}
	//instanceIds = append(instanceIds, -1)
	tasks, err = models.GetTasksBetween(instanceIds, start, end)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetDeploy(hostIp string) (task models.Task, err error) {
	waitingTasks, err := models.GetWaitingDeployTasks()
	if err != nil {
		return task, err
	}

	instances, err := models.GetInstancesByIp(hostIp)
	if err != nil {
		return task, err
	}

	for _, task := range waitingTasks {
		for _, instance := range instances {
			if task.ReplicaId == instance.ID {
				task.Instance = instance
				task.TaskStatus = define.TaskRunning
				err = models.UpdateTask(task)
				if err != nil {
					return task, err
				}
				return task, nil
			}
		}
	}
	return task, nil
}

func UpdateTask(task models.Task) error {
	if err := models.UpdateTask(task); err != nil {
		return err
	}
	return nil
}

func ADD(task models.Task) error {
	//task := &Task{
	//	ReplicaId:       data.ReplicaId,
	//	TaskType:        data.TaskType,
	//	TaskDetails:     data.TaskDetails,
	//	TaskStatus:      data.TaskStatus,
	//	ExecutionResult: data.ExecutionResult,
	//	Retry:           data.Retry,
	//	CreateAt:      data.CreateAt,
	//	UpdatedAt:      data.UpdatedAt,
	//}
	if err := models.AddTask(task); err != nil {
		return err
	}

	return nil
}

func DeleteTaskById(id int64) error {
	if err := models.DeleteTaskById(id); err != nil {
		return err
	}
	return nil
}
