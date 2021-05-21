package instance_service

import (
	"easycache/models"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/pkg/setting"
	"easycache/pkg/util"
	"easycache/service/config_service"
	"easycache/service/resource_service"
	"easycache/service/task_service"
	"easycache/service/user_service"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	jsonq "github.com/thedevsaddam/gojsonq/v2"
	"strconv"
	time2 "time"
)

func InsertInstance(instancePlus models.InstancePlus) (err error) {
	user, err := user_service.GetRedisUserById(instancePlus.Instance.RedisUserId)
	if err != nil {
		logger.Error(err)
		return err
	}
	if user.ID == 0 {
		logger.Error(err)
		return errors.New("cannot find redis user")
	}
	// todo 去除tips
	instancePlus.Uuid = user.Uuid
	if instancePlus.Instance.IsCluster == "" {
		instancePlus.Instance.IsCluster = "no"
	}
	if instancePlus.Instance.IsVirtual == "" {
		instancePlus.Instance.IsVirtual = "false"
	}
	if instancePlus.Instance.Name == "" {
		instancePlus.Instance.Name = instancePlus.Instance.IP + ":" + instancePlus.Instance.Port
	}
	if instancePlus.Instance.Status == "" {
		instancePlus.Instance.Status = define.TaskRunning
	}
	err = models.InsertInstance(instancePlus.Instance)
	if err != nil {
		logger.Error(err)
		return err
	}
	instanceId := instancePlus.Instance.ID
	instancePlus.Config.InstanceId = instanceId
	if instancePlus.Config.Content == "" {
		instancePlus.Config.Content = "{}"
	}
	err = models.InsertInstanceConfig(instancePlus.Config)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = doDeploy(instancePlus)
	if err != nil {
		logger.Error(err)
		return err
	}
	var content int64
	maxMemory := jsonq.New().FromString(instancePlus.Config.Content).Find("maxmemory.value")
	if maxMemoryStr := util.Interface2String(maxMemory); maxMemoryStr != "" {
		content, err = strconv.ParseInt(maxMemoryStr, 10, 64)
		if err != nil {
			logger.Error("InsertInstance parseInt err", err)
			return err
		}
	}
	resourceHandle := models.ResourceHandle{
		Ip:      instancePlus.Instance.IP,
		Content: content,
		Action:  "add",
	}
	err = resource_service.UpdateAllocatedMemory(resourceHandle)
	if err != nil {
		logger.Error("InsertInstance UpdateAllocatedMemory err:", err)
		return err
	}
	// todo 更新port
	return nil
}

func InsertInstanceCluster(instancePlus models.InstancePlus) (err error) {
	user, err := user_service.GetRedisUserById(instancePlus.Instance.RedisUserId)
	if err != nil {
		logger.Error(err)
		return err
	}
	if user.ID == 0 {
		logger.Error(err)
		return errors.New("cannot find redis user")
	}
	// todo 去除tips
	instancePlus.Uuid = user.Uuid
	if instancePlus.Instance.IsCluster == "" {
		instancePlus.Instance.IsCluster = "no"
	}
	if instancePlus.Instance.IsVirtual == "" {
		instancePlus.Instance.IsVirtual = "false"
	}
	if instancePlus.Instance.Name == "" {
		instancePlus.Instance.Name = instancePlus.Instance.IP + ":" + instancePlus.Instance.Port
	}
	if instancePlus.Instance.Status == "" {
		instancePlus.Instance.Status = define.TaskRunning
	}
	err = models.InsertInstanceAndUpdate(instancePlus.Instance)
	if err != nil {
		logger.Error(err)
		return err
	}
	instanceId := instancePlus.Instance.ID
	instancePlus.Config.InstanceId = instanceId
	if instancePlus.Config.Content == "" {
		instancePlus.Config.Content = "{}"
	}
	// 解base64
	bytesContent, err := base64.StdEncoding.DecodeString(instancePlus.Config.Content)
	if err != nil {
		logger.Error(err)
		return err
	}
	instancePlus.Config.Content = string(bytesContent)
	logger.Info("instancePlus.Config.Content", instancePlus.Config.Content)
	err = models.InsertInstanceConfig(instancePlus.Config)
	if err != nil {
		logger.Error(err)
		return err
	}
	err = doDeploy(instancePlus)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func doDeploy(instancePlus models.InstancePlus) (err error) {
	config := make(map[string]interface{})
	server := make(map[string]interface{})
	agentConf := make(map[string]interface{})

	agentConf["instance_id"] = instancePlus.Instance.ID
	agentConf["instance_name"] = instancePlus.Instance.Name
	agentConf["redis_user_id"] = instancePlus.Instance.RedisUserId
	agentConf["uuid"] = instancePlus.Uuid
	port, err := strconv.ParseInt(instancePlus.Instance.Port, 10, 64)
	if err != nil {
		logger.Error("doDeploy parseInt err:", err)
		return err
	}
	agentConf["redis_port"] = port
	agentConf["redis_ip"] = instancePlus.Instance.IP
	agentConf["redis_cluster_enabled"] = instancePlus.Instance.IsCluster
	agentConf["master_host"] = setting.AppSetting.PrefixUrl
	agentConf["machine_type"] = instancePlus.Instance.IsVirtual == "true"
	agentConf["redis_password"] = instancePlus.Instance.Password

	config["server"] = server
	config["agent_conf"] = agentConf
	// 防止Marshal的时候转义
	config["redis_conf"] = json.RawMessage(instancePlus.Config.Content)
	logger.Debug("service_instance doDeploy instancePlus.Config.Content", instancePlus.Config.Content)
	confJsonBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	encodeStr := base64.StdEncoding.EncodeToString(confJsonBytes)
	logger.Debug(encodeStr)
	var agentCommand string
	// 是否是docker启动
	if instancePlus.Instance.IsVirtual == "true" {
		clusterPort := port + 10000
		imageName := setting.AppSetting.ImageName
		logger.Debug("instance_service doDeploy", imageName)
		cpu := int64(1000000)
		memory := int64(274877906944)
		if cpuBar := instancePlus.Resource.Cpu; cpuBar != 0 {
			cpu = cpuBar * 1000
		}
		if memBar := instancePlus.Resource.Memory; memBar != 0 {
			memory = memBar
		}
		memStr := fmt.Sprintf("%dB", memory)
		dir := "/usr/local/services/redis"
		agentBin := "/usr/local/bin/agent"
		agentMode := "deployAndStart"
		netMode := "host"
		time := time2.Now().Format("15-04-02")
		container := fmt.Sprintf("redis_%d_%s", port, time)
		// todo 判断--user redis 是否有必要
		agentCommand = fmt.Sprintf("docker run -t -i -d --privileged=true -m %s --cpu-period 1000000 --cpu-quota %d --memory-swap -1 --name=%s -p %d:%d -p %d:%d --net=%s -v %s:%s -v /data:/data %s \"%s %s %s &\"",
			memStr, cpu, container, port, port, clusterPort,
			clusterPort, netMode, dir, dir, imageName, agentBin,
			agentMode, encodeStr)
		logger.Debug(agentCommand)
	} else {
		agentCommand = fmt.Sprintf("nohup /usr/local/bin/agent deployAndStart %s >> /usr/local/services/redis/%d.log 2>&1 &", encodeStr, port)
	}
	task := models.Task{
		ReplicaId:       instancePlus.Instance.ID,
		TaskDetails:     agentCommand,
		ExecutionResult: "init",
		TaskType:        define.TaskDeploy,
	}
	if instancePlus.IsAutoDeploy == "yes" {
		task.TaskStatus = define.TaskWaiting
	} else {
		task.TaskStatus = define.TaskFinished
	}
	err = task_service.ADD(task)
	if err != nil {
		return err
	}
	return nil
}

func GetActiveInstancesByIp(ip string) (instances []models.Instance, err error) {
	instances, err = models.GetActiveInstancesByIp(ip)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func GetInstanceById(id int64) (instance models.Instance, err error) {
	instance, err = models.GetInstanceById(id)
	if err != nil {
		logger.Error("instance_service err: ", err)
		return instance, err
	}
	return instance, nil
}

func GetInstances() (instances []models.Instance, err error) {
	instances, err = models.GetInstances()
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func GetActiveInstancesBySocket(ip string, port int) (instances []models.Instance, err error) {
	instances, err = models.GetActiveInstancesBySocket(ip, port)
	if err != nil {
		return nil, err
	}
	return instances, nil
}

func StartInstance(instance models.Instance) (err error) {
	taskData := make(map[string]interface{})
	taskData["state"] = "start"
	jsonBytes, err := json.Marshal(taskData)
	task := models.Task{ReplicaId: instance.ID, TaskDetails: string(jsonBytes), TaskStatus: define.TaskWaiting, TaskType: define.TaskChangeState}
	err = task_service.ADD(task)
	instance.Status = define.TaskRunning
	//变更任务状态
	err = UpdateInstanceStatus(instance)
	if err != nil {
		logger.Error("StartInstance UpdateInstanceStatus err:", err)
		return err
	}
	return nil
}

func StopInstance(instance models.Instance) (err error) {
	taskData := make(map[string]interface{})
	taskData["state"] = "shutdown"
	jsonBytes, err := json.Marshal(taskData)
	task := models.Task{ReplicaId: instance.ID, TaskDetails: string(jsonBytes), TaskStatus: define.TaskWaiting, TaskType: define.TaskChangeState}
	err = task_service.ADD(task)
	instance.Status = define.TaskStop
	err = UpdateInstanceStatus(instance)
	if err != nil {
		logger.Error("StartInstance UpdateInstanceStatus err:", err)
		return err
	}
	return nil
}

func DeleteInstance(instance models.Instance) (err error) {
	taskData := make(map[string]interface{})
	taskData["state"] = "delete"
	jsonBytes, err := json.Marshal(taskData)
	task := models.Task{ReplicaId: instance.ID, TaskDetails: string(jsonBytes), TaskStatus: define.TaskWaiting, TaskType: define.TaskChangeState}
	err = task_service.ADD(task)
	// 更新资源
	config, err := config_service.GetConfigByReplicaId(instance.ID)
	if err != nil {
		logger.Error("DeleteInstance GetConfigByReplicaId err:", err)
		return err
	}
	var content int64
	maxMemory := jsonq.New().FromString(config.Content).Find("maxmemory.value")
	if maxMemoryStr := util.Interface2String(maxMemory); maxMemoryStr != "" {
		content, err = strconv.ParseInt(maxMemoryStr, 10, 64)
		if err != nil {
			logger.Error("DeleteInstance parseInt err", err)
			return err
		}
	}
	resourceHandle := models.ResourceHandle{
		Ip:      instance.IP,
		Content: content,
		Action:  "delete",
	}
	err = resource_service.UpdateAllocatedMemory(resourceHandle)
	if err != nil {
		logger.Error("DeleteInstance UpdateAllocatedMemory err:", err)
		return err
	}
	instance.Status = define.TaskDelete
	//变更任务状态
	err = UpdateInstanceStatus(instance)
	if err != nil {
		logger.Error("DeleteInstance UpdateInstanceStatus err:", err)
		return err
	}
	return nil
}

func UpdateInstanceStatus(instance models.Instance) (err error) {
	err = models.UpdateTaskStatus(instance)
	if err != nil {
		logger.Error("UpdateInstanceStatus err:", err)
		return err
	}
	return nil
}
