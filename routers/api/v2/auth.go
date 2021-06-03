package v2

import (
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/auth_service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(context *gin.Context) {
	appG := app.Gin{C: context}
	username := context.Query("username")
	password := context.Query("password")

	flag, err := auth_service.Login(username, password)
	if err != nil {
		logger.Error("Login error", err)
		appG.Response(http.StatusOK, define.ErrorLogin, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, flag)
}

func Register(context *gin.Context) {
	appG := app.Gin{C: context}
	username := context.Query("username")
	password := context.Query("password")

	flag, err := auth_service.Register(username, password)
	if err != nil {
		logger.Error("Register error", err)
		appG.Response(http.StatusOK, define.ErrorRegister, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, flag)
}

func GetUser(context *gin.Context) {
	appG := app.Gin{C: context}
	username := context.Query("username")

	auth, err := auth_service.GetUser(username)
	if err != nil {
		logger.Error("Login error", err)
		appG.Response(http.StatusOK, define.ErrorLogin, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, auth)
}
