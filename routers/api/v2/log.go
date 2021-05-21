package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/log_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertLogs(context *gin.Context) {
	app2 := app.Gin{C: context}
	var log models.LogInfo
	if err := context.BindJSON(&log); err != nil {
		logger.Error("router InsertLogs err:", err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Debug("RedisUserId: ", log.RedisUserId, "IP", log.InstanceIP, "Port", log.InstancePort)
	err := log_service.InsertLogs(log)
	if err != nil {
		logger.Error("router InsertLogs", err)
		app2.Response(http.StatusOK, define.ErrorInsertLog, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func GetLogs(c *gin.Context) {
	appG := app.Gin{C: c}
	ip := c.Query("ip")
	port := c.Query("port")
	logger.Info("ip:", ip, "port", port)
	infos, err := log_service.GetAllLogs(ip, port, 20)
	if err != nil {
		appG.Response(http.StatusOK, define.ErrorGetLogs, nil)
		return
	}

	appG.Response(http.StatusOK, define.SUCCESS, infos)
}

func DeleteLogs(c *gin.Context) {
	appG := app.Gin{C: c}
	ip := c.Query("ip")
	port := c.Query("port")
	logger.Info("ip:", ip, "port", port)
	err := log_service.DeleteLogs(ip, port)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorDeleteLogs, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, nil)
	return
}
