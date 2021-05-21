package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/instance_service"
	"easycache/service/user_service"
	"encoding/json"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
)

// InsertInstance ...
// redis_user_id
// instance
// config 除了md5
func InsertInstance(c *gin.Context) {
	app2 := app.Gin{C: c}
	var instancePlus models.InstancePlus
	if err := c.ShouldBind(&instancePlus); err != nil {
		logger.Error(err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	redisUserId := com.StrTo(c.Param("redis_user_id")).MustInt64()
	logger.Debug("username: ", instancePlus.UserName, "redis_user_id", redisUserId, "name", instancePlus.Instance.Name)
	user, err := user_service.GetRedisUserById(redisUserId)
	if err != nil || user.ID == 0 {
		logger.Error(err)
		app2.Response(http.StatusOK, define.ErrorGetRedisUser, "GetRedisUserById error")
		return
	}
	err = instance_service.InsertInstance(instancePlus)
	if err != nil {
		logger.Error(err)
		app2.Response(http.StatusOK, define.ErrorInsertInstance, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func GetInstances(c *gin.Context) {
	appG := app.Gin{C: c}
	instance, err := instance_service.GetInstances()
	if err != nil {
		appG.Response(http.StatusInternalServerError, define.ErrorGetInstance, nil)
		return
	}
	rspJson, err := json.Marshal(instance)
	if err != nil {
		appG.Response(http.StatusInternalServerError, define.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, string(rspJson))
}

func StartInstance(c *gin.Context) {
	appG := app.Gin{C: c}
	logger.Info(c.Request)
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, nil)
		return
	}
	logger.Info("id", id)
	instance, err := instance_service.GetInstanceById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetInstance, err.Error())
		return
	}

	err = instance_service.StartInstance(instance)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorStartInstance, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, instance)
	return
}

func StopInstance(c *gin.Context) {
	appG := app.Gin{C: c}
	logger.Info(c.Request)
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, nil)
		return
	}
	logger.Info("id", id)
	instance, err := instance_service.GetInstanceById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetInstance, err.Error())
		return
	}

	err = instance_service.StopInstance(instance)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorStopInstance, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, instance)
	return
}

func DeleteInstance(c *gin.Context) {
	appG := app.Gin{C: c}
	logger.Info(c.Request)
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, define.INVALID_PARAMS, nil)
		return
	}
	logger.Info("id", id)
	instance, err := instance_service.GetInstanceById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetInstance, err.Error())
		return
	}

	err = instance_service.DeleteInstance(instance)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorDeleteInstance, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, instance)
	return
}
