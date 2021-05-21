package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/alarm_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertAlarmInfo(context *gin.Context) {
	app2 := app.Gin{C: context}
	var alarm models.AlarmInfo
	if err := context.BindJSON(&alarm); err != nil {
		logger.Error("router InsertAlarmInfo err:", err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Debug("RedisUserId: ", alarm.RedisUserId, "IP", alarm.IP, "Port", alarm.Port)
	err := alarm_service.InsertAlarmInfo(alarm)
	if err != nil {
		logger.Error("router InsertAlarmInfo", err)
		app2.Response(http.StatusOK, define.ErrorInsertAlarm, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func GetAlarmInfos(c *gin.Context) {
	appG := app.Gin{C: c}
	infos, err := alarm_service.GetAllAlarmInfos()
	if err != nil {
		appG.Response(http.StatusOK, define.ErrorGetAlarms, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, infos)
}

func DeleteAlarmInfos(c *gin.Context) {
	appG := app.Gin{C: c}
	err := alarm_service.DeleteAlarmInfos()
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorDeleteAlarms, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, nil)
	return
}
