package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/resource_service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InsertResource(context *gin.Context) {
	app2 := app.Gin{C: context}
	var resource models.ResourceInfo
	if err := context.ShouldBind(&resource); err != nil {
		logger.Error("router InsertResource err:", err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Debug("HostIp: ", resource.HostIp, "LimitedMemory", resource.LimitedMemory, "TotalHostMemory", resource.TotalHostMemory)
	err := resource_service.InsertResource(resource)
	if err != nil {
		logger.Error("router_resource InsertResource", err)
		app2.Response(http.StatusOK, define.ErrorInsertResource, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func GetResources(c *gin.Context) {
	appG := app.Gin{C: c}
	resources, err := resource_service.GetAllResources()
	if err != nil {
		appG.Response(http.StatusOK, define.ErrorAllResource, nil)
		return
	}
	rspJson, err := json.Marshal(resources)
	if err != nil {
		appG.Response(http.StatusOK, define.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, string(rspJson))
}

func DeleteResource(c *gin.Context) {
	appG := app.Gin{C: c}
	logger.Info(c.Request)
	ip := c.Query("host_ip")
	logger.Info("ip", ip)
	_, err := resource_service.GetResourceInfoByHost(ip)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetResource, err.Error())
		return
	}

	err = resource_service.DeleteResource(ip)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorDeleteResource, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, nil)
	return
}
