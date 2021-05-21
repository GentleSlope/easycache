package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/monitor_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func InsertMonitorInfo(c *gin.Context) {
	app2 := app.Gin{C: c}
	var info models.MonitorInfo
	if err := c.BindJSON(&info); err != nil {
		logger.Error("router InsertMonitorInfo err:", err)
		app2.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}
	logger.Debug("replica_ip: ", info.ReplicaIp, "redis_user_id", info.RedisUserId, "info", info.Info)
	err := monitor_service.InsertInfo(info)
	if err != nil {
		logger.Error("router_resource InsertResource", err)
		app2.Response(http.StatusInternalServerError, define.ErrorInsertInfo, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func GetMonitorInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	ip := c.Query("ip")
	//redisUserId, err := strconv.ParseInt(c.Query("redis_user_id"), 10, 64)
	//if err != nil {
	//	logger.Error("GetMonitorInfo ParseInt err:", err)
	//	appG.Response(http.StatusOK, define.ERROR, err.Error())
	//	return
	//}
	port, err := strconv.ParseInt(c.Query("port"), 10, 64)
	if err != nil {
		logger.Error("GetMonitorInfo ParseInt err:", err)
		appG.Response(http.StatusOK, define.ERROR, err.Error())
		return
	}
	logger.Info("ip:", ip, "port", port)
	infos, err := monitor_service.GetMonitorInfo(ip, port, 10)
	if err != nil {
		appG.Response(http.StatusOK, define.ErrorGetInstance, nil)
		return
	}

	appG.Response(http.StatusOK, define.SUCCESS, infos)
}
