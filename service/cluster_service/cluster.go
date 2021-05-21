package cluster_service

import (
	"easycache/models"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/instance_service"
	"easycache/service/task_service"
	"easycache/service/user_service"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"net"
	"os/exec"
	"sync"
	"time"
)

func GetAllClusters() (clusters []models.Cluster, err error) {
	clusters, err = models.GetAllClusters()
	if err != nil {
		return clusters, err
	}
	return clusters, nil
}

func GetClusterInstancesById(id int64) (instances []models.InstanceDisplay, err error) {
	// 首先查询集群节点
	clusterInstances, err := models.GetClusterInstancesById(id)
	if err != nil {
		return instances, err
	}
	for _, clusterInstance := range clusterInstances {
		instance, err := models.GetInstanceById(clusterInstance.InstanceId)
		if err != nil {
			return instances, err
		}
		instanceDisplay := models.InstanceDisplay{
			Role:     clusterInstance.Role,
			MasterId: clusterInstance.MasterId,
			Instance: instance,
		}
		instances = append(instances, instanceDisplay)
	}
	return instances, nil
}

func DeleteCluster(id int64) error {
	if err := models.DeleteCluster(id); err != nil {
		return err
	}
	return nil
}

func AddSlaveNode(masterId int64, instanceId int64, clusterId int64) (err error) {
	master, err := instance_service.GetInstanceById(masterId)
	if err != nil {
		return err
	}
	slave, err := instance_service.GetInstanceById(instanceId)
	if err != nil {
		return err
	}
	if define.TaskRunning != master.Status || define.TaskRunning != slave.Status {
		return errors.New("node is not running!")
	}
	masterSocket := fmt.Sprintf("%s:%s", master.IP, master.Port)
	slaveSocket := fmt.Sprintf("%s:%s", slave.IP, slave.Port)
	config := make(map[string]string)
	config["action"] = "addSlave"
	config["add_endpoint"] = slaveSocket
	config["master_endpoint"] = masterSocket
	config["password"] = slave.Password

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		logger.Error("cluster_svr AddSlaveNode Marshal nodes err:", err)
		return err
	}
	logger.Info(string(jsonBytes))
	encodeStr := base64.StdEncoding.EncodeToString(jsonBytes)
	command := "/usr/local/bin/cluster-op"
	logger.Info(fmt.Sprintf("%s %s", command, encodeStr))
	err = exec.Command(command, encodeStr).Run()
	if err != nil {
		logger.Error("cluster_svr AddSlaveNode exec command err:", err)
		return err
	}
	err = models.InsertClusterInstance(models.ClusterInstance{
		ClusterId:  clusterId,
		InstanceId: slave.ID,
		Role:       "slave",
		MasterId:   masterId,
	})
	if err != nil {
		logger.Error("cluster_svr InsertClusterInstance exec command err:", err)
		return err
	}
	return nil
}

func DeleteSlaveNode(instanceId int64) (err error) {
	slave, err := instance_service.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	clusterInstance, err := models.GetClusterInstancesByInstanceId(instanceId)
	if err != nil {
		return err
	}
	masterId := clusterInstance.MasterId
	master, err := instance_service.GetInstanceById(masterId)
	if err != nil {
		return err
	}

	if define.TaskRunning != master.Status || define.TaskRunning != slave.Status {
		return errors.New("node is not running!")
	}
	masterSocket := fmt.Sprintf("%s:%s", master.IP, master.Port)
	slaveSocket := fmt.Sprintf("%s:%s", slave.IP, slave.Port)
	config := make(map[string]string)
	config["action"] = "forget"
	config["forget_endpoint"] = slaveSocket
	config["cluster_endpoint"] = masterSocket
	config["password"] = slave.Password

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		logger.Error("cluster_svr AddSlaveNode Marshal nodes err:", err)
		return err
	}
	logger.Info(string(jsonBytes))
	encodeStr := base64.StdEncoding.EncodeToString(jsonBytes)
	command := "/usr/local/bin/cluster-op"
	logger.Info(fmt.Sprintf("%s %s", command, encodeStr))
	err = exec.Command(command, encodeStr).Run()
	if err != nil {
		logger.Error("cluster_svr AddSlaveNode exec command err:", err)
		return err
	}

	// 级联删除
	err = models.DeleteClusterInstance(models.ClusterInstance{InstanceId: instanceId})
	if err != nil {
		logger.Error("cluster_svr DeleteClusterInstance err:", err)
		return err
	}

	err = instance_service.DeleteInstance(slave)
	if err != nil {
		logger.Error("cluster_svr DeleteInstance err:", err)
		return err
	}
	return nil
}

func AddMasterNode(instanceId int64, clusterId int64) (err error) {
	master, err := instance_service.GetInstanceById(instanceId)
	if err != nil {
		return err
	}
	masterSocket := fmt.Sprintf("%s:%s", master.IP, master.Port)
	instances, err := GetClusterInstancesById(clusterId)
	if err != nil {
		return err
	}
	config := make(map[string]interface{})
	var sockets []map[string]string
	var pwd string
	for _, instance := range instances {
		slaveConf := make(map[string]string)
		slaveSocket := fmt.Sprintf("%s:%s", instance.IP, instance.Port)
		slaveConf["endpoint"] = slaveSocket
		sockets = append(sockets, slaveConf)
		pwd = instance.Password
	}
	config["action"] = "addMaster"
	config["add_endpoint"] = masterSocket
	config["password"] = pwd
	config["cluster_masters"] = sockets

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		logger.Error("cluster_svr AddMasterNode Marshal nodes err:", err)
		return err
	}
	logger.Info(string(jsonBytes))
	encodeStr := base64.StdEncoding.EncodeToString(jsonBytes)
	command := "/usr/local/bin/cluster-op"
	logger.Info(fmt.Sprintf("%s %s", command, encodeStr))
	err = exec.Command(command, encodeStr).Run()
	if err != nil {
		logger.Error("cluster_svr AddMasterNode exec command err:", err)
		return err
	}
	err = models.InsertClusterInstance(models.ClusterInstance{
		ClusterId:  clusterId,
		InstanceId: master.ID,
		Role:       "master",
		MasterId:   master.ID,
	})
	if err != nil {
		logger.Error("cluster_svr InsertClusterInstance exec command err:", err)
		return err
	}
	return nil
}

func CreateCluster(cluster models.ClusterForm) (err error) {
	//signature := true
	redisUserId := cluster.RedisUserId
	// 官方推荐使用三主三从
	if len(cluster.Instances) != 6 {
		return errors.New("instances size err")
	}
	if len(cluster.Password) < 6 {
		return errors.New("password size err")
	}
	// 检查是否已经被占用
	for _, instance := range cluster.Instances {
		ip := instance.Ip
		port := instance.Port
		instances, err := instance_service.GetActiveInstancesBySocket(ip, port)
		if err != nil {
			logger.Error("cluster_service CreateCluster GetActiveInstancesBySocket err: ", err)
			return err
		}
		if len(instances) > 0 {
			errMag := "cluster_service CreateCluster GetActiveInstancesBySocket err: ip: port is used"
			logger.Error(errMag)
			return errors.New(errMag)
		}
	}
	// 看是否能ping通客户端
	for _, instance := range cluster.Instances {
		conn, err := net.DialTimeout("tcp", instance.Ip+":22", 3*time.Second)
		if err != nil {
			logger.Error("cluster_service CreateCluster ping ip", instance.Ip, "err: ", err)
			return err
		}
		_ = conn.Close()
	}
	baseVersion := cluster.Instances[0].Version
	for _, instance := range cluster.Instances {
		if instance.Version != baseVersion {
			errMag := "cluster_service CreateCluster Version Check err: not assign"
			logger.Error(errMag)
			return errors.New(errMag)
		}
	}
	instanceConfig := models.NodeConfig{Content: cluster.InstanceConfigExtension.Config.Content}
	if cluster.RedisUserId == 0 {
		redisUser := models.RedisUser{Name: cluster.RedisUserName}
		err = user_service.InsertUser(&redisUser)
		if err != nil {
			logger.Error("cluster_service CreateCluster InsertUser err", err)
			return err
		}
		redisUserId = redisUser.ID
		logger.Info("cluster_service CreateCluster RedisUserId:", redisUserId)
	}
	// process instance
	instancePlus := models.InstancePlus{
		IsAutoDeploy: "yes",
		Instance: models.Instance{
			IsCluster:   "yes",
			IsVirtual:   "true",
			Password:    cluster.Password,
			RedisUserId: redisUserId,
		},
		Config: instanceConfig,
	}
	ids := make(map[int64]string)
	for _, v := range cluster.Instances {
		instancePlus.Instance.ID = v.ID
		instancePlus.Instance.IP = v.Ip
		instancePlus.Instance.Port = fmt.Sprintf("%d", v.Port)
		instancePlus.Instance.Name = fmt.Sprintf("%s:%d", v.Ip, v.Port)
		instancePlus.Instance.Version = v.Version
		// todo 这种插入有并发问题
		err := instance_service.InsertInstanceCluster(instancePlus)
		if err != nil {
			logger.Error("cluster_service CreateCluster InsertInstanceCluster err", err)
			return err
		}
		logger.Info(instancePlus.Instance.ID)
		ids[instancePlus.Instance.ID] = fmt.Sprintf("%s:%d", v.Ip, v.Port)
	}
	waitGroup := sync.WaitGroup{}
	var errInternal error
	ok := true
	// 重试最多100次
	for i := 0; i < 100; i++ {
		time.Sleep(2 * time.Second)
		ok = true
		for k := range ids {
			waitGroup.Add(1)
			go func() {
				defer func() {
					if e := recover(); e != nil {
						logger.Error("deploy task error", e)
						errInternal = errors.New("deploy task error")
						waitGroup.Done()
					}
				}()
				task, err := task_service.GetTaskByReplicaId(k)
				if err != nil {
					panic(err)
				}
				if task.TaskStatus == define.TaskFailed {
					panic(errors.New(fmt.Sprintf("instance: %d deploy failed", k)))
				}
				// todo 修改状态名称
				if task.TaskStatus != define.TaskFinished {
					ok = false
					logger.Error(errors.New(fmt.Sprintf("instance: %d is deploying", k)))
				}
				waitGroup.Done()
			}()
		}
		waitGroup.Wait()
		if errInternal != nil {
			return errInternal
		}
		if ok {
			logger.Info("all tasks down")
			break
		}
	}
	for i := 0; i < 100; i++ {
		time.Sleep(20 * time.Second)
		pwd := cluster.Password
		logger.Info("pwd", pwd)
		for _, v := range cluster.Instances {
			err = connect(fmt.Sprintf("%s:%d", v.Ip, v.Port), pwd)
			if err != nil {
				logger.Error(fmt.Sprintf("%s:%d", v.Ip, v.Port), " instance connect err: ", err)
				ok = false
				break
			}
		}
		if ok {
			logger.Info("all redis instances check ok")
			break
		}
	}
	if !ok {
		return errors.New("redis instances not ready")
	}
	clusterDeploy := models.CLusterDeploy{
		ClusterPwd: cluster.Password,
		Cluster: models.Cluster{
			Name:        cluster.ClusterName,
			SlotMode:    cluster.SlotMode,
			MasterNum:   6,
			RedisUserId: redisUserId,
		},
	}
	var clusterInstances []models.ClusterInstance
	for _, instance := range cluster.Instances {
		logger.Info("instances:", instance)
		node := models.ClusterInstance{}
		for k, v := range ids {
			if v == fmt.Sprintf("%s:%d", instance.Ip, instance.Port) {
				node.InstanceId = k
			}
			if instance.MasterIp != "" && v == fmt.Sprintf("%s:%d", instance.MasterIp, instance.MasterPort) {
				node.MasterId = k
			}
		}
		if instance.Role == "slave" {
			node.Role = "slave"
		} else {
			node.Role = "master"
		}
		clusterInstances = append(clusterInstances, node)
	}
	clusterDeploy.Instances = clusterInstances
	err = insertCluster(clusterDeploy)
	if err != nil {
		logger.Error(err)
		return err
	}
	return
}

func insertCluster(clusterDeploy models.CLusterDeploy) (err error) {
	var nodesNum, masterNum int64
	for _, v := range clusterDeploy.Instances {
		if v.Role == "master" {
			nodesNum++
			masterNum++
		}
		if v.Role == "slave" {
			nodesNum++
		}
	}
	clusterDeploy.Cluster.MasterNum = masterNum
	clusterDeploy.Cluster.NodesNum = nodesNum
	err = models.InsertCluster(&clusterDeploy.Cluster)
	if err != nil {
		logger.Error("cluster_svr insertCluster err:", err)
		return err
	}
	// todo 版本校验
	for k, v := range clusterDeploy.Instances {
		isExist, err := models.GetClusterInstanceById(v.ClusterId, v.InstanceId)
		if err != nil {
			logger.Error("cluster_svr insertCluster GetClusterInstanceById err:", err)
			return err
		}
		if isExist {
			errMsg := fmt.Sprintf("cluster_svr insertCluster instance %v is exist in cluster %v", v.InstanceId, v.ClusterId)
			logger.Error(errMsg)
			return err
		}
		// 注意这里做了一次转换
		instance, err := models.GetInstanceById(v.InstanceId)
		if err != nil {
			logger.Error("cluster_svr insertCluster GetInstanceById err:", err)
			return err
		}
		if instance.ID == 0 {
			logger.Error("cluster_svr insertCluster GetInstanceById err")
			return errors.New("cluster_svr insertCluster GetInstanceById err")
		}
		if instance.IsCluster == "no" {
			logger.Error("cluster_svr insertCluster GetInstanceById cluster mode err")
			return errors.New("cluster_svr insertCluster GetInstanceById cluster mode err")
		}
		if k == 0 {
			clusterDeploy.ClusterPwd = instance.Password
		}
		if clusterDeploy.ClusterPwd != instance.Password {
			logger.Error("cluster_svr insertCluster GetInstanceById cluster pwd err")
			return errors.New("cluster_svr insertCluster GetInstanceById cluster pwd err")
		}
		logger.Info("clusterDeploy.Cluster:", clusterDeploy.Cluster)
		v.ClusterId = clusterDeploy.Cluster.Id
		logger.Info("clusterInstance will be inserted", v)
		err = models.InsertClusterInstance(v)
		if err != nil {
			logger.Error("cluster_svr insertCluster InsertClusterInstance err:", err)
			return err
		}
	}
	err = deployCluster(clusterDeploy)
	if err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func deployCluster(clusterDeploy models.CLusterDeploy) (err error) {
	config := make(map[string]interface{})
	clusterInfo := make(map[string]interface{})
	var nodes []map[string]interface{}

	config["clusterInfo"] = clusterInfo
	config["action"] = "create"

	clusterInfo["slot_mode"] = clusterDeploy.Cluster.SlotMode
	clusterInfo["nodes_num"] = clusterDeploy.Cluster.NodesNum
	clusterInfo["master_num"] = clusterDeploy.Cluster.MasterNum
	clusterInfo["cluster_password"] = clusterDeploy.ClusterPwd

	for _, clusterInstance := range clusterDeploy.Instances {
		node := make(map[string]interface{})
		nodeInstance, err := models.GetInstanceById(clusterInstance.InstanceId)
		if err != nil {
			logger.Error("cluster_svr insertCluster deployCluster GetInstanceById err:", err)
			return err
		}
		node["endpoint"] = fmt.Sprintf("%v:%v", nodeInstance.IP, nodeInstance.Port)
		node["slots"] = clusterInstance.Slots
		node["role"] = clusterInstance.Role
		if clusterInstance.Role == "slave" {
			master, err := models.GetInstanceById(clusterInstance.MasterId)
			if err != nil {
				logger.Error("cluster_svr insertCluster deployCluster GetInstanceById master err:", err)
				return err
			}
			node["my_master_endpoint"] = fmt.Sprintf("%v:%v", master.IP, master.Port)
		}
		nodes = append(nodes, node)
	}
	config["nodes"] = nodes
	jsonBytes, err := json.Marshal(config)
	if err != nil {
		logger.Error("cluster_svr insertCluster deployCluster Marshal nodes err:", err)
		return err
	}
	logger.Info(string(jsonBytes))
	encodeStr := base64.StdEncoding.EncodeToString(jsonBytes)
	command := "/usr/local/bin/cluster-op"
	logger.Info(fmt.Sprintf("%s %s", command, encodeStr))
	err = exec.Command(command, encodeStr).Run()
	if err != nil {
		logger.Error("cluster_svr insertCluster deployCluster exec command err:", err)
		return err
	}
	return nil
}

// 初始化连接
func connect(add string, password string) error {
	logger.Info("address: ", add)
	logger.Info("password: ", password)
	rdb := redis.NewClient(&redis.Options{
		Addr:        add,
		Password:    password,
		DB:          0, // use default DB
		DialTimeout: 2 * time.Second,
	})
	defer rdb.Close()
	pong, err := rdb.Ping().Result()
	if err != nil {
		logger.Error("connect err", pong, err)
		return err
	}
	return nil
}
