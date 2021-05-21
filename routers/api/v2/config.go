package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/config_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetNodeConfigByReplicaId(c *gin.Context) {
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

	config, err := config_service.GetConfigByReplicaId(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetConfig, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, config)
}
func UpdateNodeConfig(c *gin.Context) {
	appG := app.Gin{C: c}
	var config models.NodeConfig
	err := c.BindJSON(&config)
	if err != nil {
		appG.Response(http.StatusOK, define.INVALID_PARAMS, err)
		return
	}
	instanceId := config.InstanceId
	content := config.Content
	logger.Info("instance_id", instanceId, "content", content)
	currentConfig, err := config_service.GetConfigByReplicaId(instanceId)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetConfig, err.Error())
		return
	}
	if currentConfig.InstanceId == 0 {
		logger.Error("args error!")
		appG.Response(http.StatusOK, define.ErrorGetConfig, "args error!")
		return
	}
	currentConfig.Content = content
	err = config_service.UpdateConfig(currentConfig)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorUpdateConfig, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, "ok")
	return
}
