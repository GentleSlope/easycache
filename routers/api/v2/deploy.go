package v2

import (
	"easycache/models"
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/task_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)

type TaskResponse struct {
	TaskId          int64  `form:"task_id" valid:"Required;Min(1)"`
	TaskStatus      string `form:"task_status" valid:"Required;MaxSize(100)"`
	ExecutionResult string `form:"execution_result" valid:"Required;MaxSize(255)"`
}

func GetDeployTask(c *gin.Context) {
	appG := app.Gin{C: c}
	uuid := com.StrTo(c.Param("uuid")).String()
	ip := com.StrTo(c.Param("ip")).String()

	valid := validation.Validation{}
	valid.IP(ip, "hostIP")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}

	if uuid != define.ServerUuId {
		logger.Error("UUid it not match!")
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
	}

	task, err := task_service.GetDeploy(ip)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorGetTasksFail, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, task)
}

func UpdateDeployTask(c *gin.Context) {
	appG := app.Gin{C: c}

	form := TaskResponse{}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != define.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	status := form.TaskStatus
	taskId := form.TaskId
	currentTask, err := task_service.GetTaskById(taskId)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorGetTasksFail, nil)
		return
	}
	if currentTask.ID == 0 || define.TaskDeploy != currentTask.TaskType {
		logger.Error("args error!")
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}
	currentTask.ID = taskId
	currentTask.TaskStatus = status
	currentTask.ExecutionResult = form.ExecutionResult

	err = task_service.UpdateTask(currentTask)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorUpdateTask, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, "ok")
}

func UpdateNotDeployTask(c *gin.Context) {
	appG := app.Gin{C: c}

	req := models.Task{}
	httpCode, errCode := app.BindAndValid(c, &req)
	if errCode != define.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	status := req.TaskStatus
	taskId := req.ID
	currentTask, err := task_service.GetTaskById(taskId)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorGetTasksFail, nil)
		return
	}
	// 这个不是一个部署任务
	if currentTask.ID == 0 || define.TaskDeploy == currentTask.TaskType {
		logger.Error("args error!")
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}
	currentTask.ID = taskId
	currentTask.TaskStatus = status
	currentTask.ExecutionResult = req.ExecutionResult

	err = task_service.UpdateTask(currentTask)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorUpdateTask, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, "ok")
}
