package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/cluster_service"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetClusters(c *gin.Context) {
	appG := app.Gin{C: c}
	clusters, err := cluster_service.GetAllClusters()
	if err != nil {
		appG.Response(http.StatusOK, define.ErrorAllCluster, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, clusters)
}

func GetClusterInstanceById(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)

	valid := validation.Validation{}
	valid.Min(id, 1, "id")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, nil)
		return
	}
	logger.Info("id", id)
	instances, err := cluster_service.GetClusterInstancesById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetClusterInstances, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, instances)
}

func DeleteCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Info("id", id)
	err = cluster_service.DeleteCluster(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorDeleteCluster, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func CreateCluster(context *gin.Context) {
	app2 := app.Gin{C: context}
	var cluster models.ClusterForm
	if err := context.BindJSON(&cluster); err != nil {
		logger.Error("router CreateCluster err:", err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Debug("RedisUserId: ", cluster.RedisUserId, "RedisUserName", cluster.RedisUserName, "ClusterName", cluster.ClusterName, "SlotMode", cluster.SlotMode)
	for _, v := range cluster.Instances {
		logger.Debug("Host: ", v.Ip+fmt.Sprintf(":%d", v.Port),
			"\nMasterHost: ", v.MasterIp+fmt.Sprintf(":%v", v.MasterPort),
			"\nRole: ", v.Role,
			"\nSlots: ", v.Slots,
			"\nVersion:", v.Version)
	}

	err := cluster_service.CreateCluster(cluster)
	if err != nil {
		logger.Error("cluster_resource CreateCluster", err)
		app2.Response(http.StatusOK, define.ErrorCreateCluster, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func ScaleCluster(context *gin.Context) {
	app2 := app.Gin{C: context}
	var masterId int64
	if context.Query("master_id") != "" {
		temp, err := strconv.ParseInt(context.Query("master_id"), 10, 64)
		masterId = temp
		if err != nil {
			logger.Error(err)
			app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
			return
		}
	}
	instanceId, err := strconv.ParseInt(context.Query("instance_id"), 10, 64)
	if err != nil {
		logger.Error(err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	clusterId, err := strconv.ParseInt(context.Query("cluster_id"), 10, 64)
	if err != nil {
		logger.Error(err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Debug("master_id: ", masterId, "instanceId", instanceId, "clusterId", clusterId)
	if masterId != 0 {
		err := cluster_service.AddSlaveNode(masterId, instanceId, clusterId)
		if err != nil {
			logger.Error(err)
			app2.Response(http.StatusOK, define.ERROR, err.Error())
			return
		}
	}
	if masterId == 0 {
		err := cluster_service.AddMasterNode(instanceId, clusterId)
		if err != nil {
			logger.Error(err)
			app2.Response(http.StatusOK, define.ERROR, err.Error())
			return
		}
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func DeleteNode(context *gin.Context) {
	app2 := app.Gin{C: context}
	instanceId, err := strconv.ParseInt(context.Query("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	err = cluster_service.DeleteSlaveNode(instanceId)
	if err != nil {
		logger.Error(err)
		app2.Response(http.StatusOK, define.ERROR, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}
