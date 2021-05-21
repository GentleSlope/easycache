package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/user_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GerRedisUsers(context *gin.Context) {
	appG := app.Gin{C: context}
	users, err := user_service.GetRedisUsers()
	if err != nil {
		appG.Response(http.StatusOK, define.ErrorAllUsers, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, users)
}

func InsertRedisUser(context *gin.Context) {
	app2 := app.Gin{C: context}
	var user models.RedisUser
	if err := context.BindJSON(&user); err != nil {
		logger.Error("router InsertRedisUser err:", err)
		app2.Response(http.StatusOK, define.INVALID_PARAMS, err.Error())
		return
	}
	logger.Debug("Name: ", user.Name, "Details", user.Details)
	err := user_service.InsertUser(&user)
	if err != nil {
		logger.Error("router InsertRedisUser", err)
		app2.Response(http.StatusOK, define.ErrorInsertUser, err.Error())
		return
	}
	app2.Response(http.StatusOK, define.SUCCESS, nil)
	return
}

func DeleteRedisUser(c *gin.Context) {
	appG := app.Gin{C: c}
	logger.Info(c.Request)
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logger.Error("router DeleteRedisUser", err)
		appG.Response(http.StatusOK, define.ERROR, err.Error())
		return
	}
	logger.Info("id", id)
	err = user_service.DeleteUserById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorInsertUser, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, nil)
	return
}
