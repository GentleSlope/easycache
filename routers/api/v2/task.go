package v2

import (
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/task_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
	"strconv"
	"time"
)

type TaskForm struct {
	//Task        models.Task `json:"task"`
	RedisUserId int    `form:"redis_user_id" valid:"Required;Min(1)"`
	Begin       string `form:"begin" valid:"Required;MaxSize(100)"`
	End         string `form:"end" valid:"Required;MaxSize(100)"`
}

func GetTaskInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	form := TaskForm{}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != define.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	start, _ := strconv.ParseInt(form.Begin, 10, 64)
	end, _ := strconv.ParseInt(form.End, 10, 64)
	logger.Debug(start, "-", end)
	//时间戳转化为日期
	startTime := time.Unix(start, 0).Format("2006-01-02 15:04:05")
	endTimeDate := time.Unix(end, 0)
	endTimeDate = endTimeDate.AddDate(0, 0, 1)
	endTime := endTimeDate.Format("2006-01-02 15:04:05")

	logger.Debug(startTime, "-", endTime)
	tasks, err := task_service.GetBetweenData(int64(int(form.RedisUserId)), startTime, endTime)

	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetTasksFail, nil)
		return
	}
	if len(tasks) == 0 {
		appG.Response(http.StatusOK, define.ErrorNotExistTask, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, tasks)
}

func GetTaskByReplicaId(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("replica_id")).MustInt64()
	begin := com.StrTo(c.Param("begin")).MustInt64()
	end := com.StrTo(c.Param("end")).MustInt64()

	startTime := com.Date(begin, "yyyy-MM-dd")
	endTime := com.Date(end, "yyyy-MM-dd")
	logger.Debug(startTime, "-", endTime)

	valid := validation.Validation{}
	valid.Min(id, 1, "replica_id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}

	tasks, err := task_service.Get(id)

	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorGetTasksFail, nil)
		return
	}
	if len(tasks) == 0 {
		appG.Response(http.StatusOK, define.ErrorNotExistTask, nil)
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, tasks)
}

func GetNodeTaskByReplicaId(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)

	valid := validation.Validation{}
	valid.Min(id, 1, "id")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, define.INVALID_PARAMS, nil)
		return
	}
	logger.Info("id", id)

	task, err := task_service.GetNodeChangeTask(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusInternalServerError, define.ErrorGetTasksFail, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, task)
}

func DeleteTask(c *gin.Context) {
	appG := app.Gin{C: c}
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetTasksFail, err.Error())
		return
	}
	logger.Info("id", id)
	_, err = task_service.GetTaskById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorGetResource, err.Error())
		return
	}

	err = task_service.DeleteTaskById(id)
	if err != nil {
		logger.Error(err)
		appG.Response(http.StatusOK, define.ErrorDeleteTask, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, nil)
	return
}
