package v2

import (
	"easycache/pkg/app"
	"easycache/pkg/define"
	"easycache/pkg/logger"
	"easycache/service/book_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func SearchDoctor(context *gin.Context) {
	appG := app.Gin{C: context}
	department := context.Query("department")
	subject := context.Query("subject")
	index, _ := strconv.ParseInt(context.Query("index"), 10, 64)

	timeNow := time.Now()
	timeBookingStart := timeNow.AddDate(0, 0, int(index))

	timeStr := timeBookingStart.Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)

	doctors, err := book_service.SearchDoctors(department, subject, t)
	if err != nil {
		logger.Error("Search doctors error", err)
		appG.Response(http.StatusOK, define.ErrorBooking, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, doctors)
}

func Booking(context *gin.Context) {
	appG := app.Gin{C: context}
	id, err := strconv.ParseInt(context.Query("doctor_id"), 10, 64)
	index, _ := strconv.ParseInt(context.Query("index"), 10, 64)

	timeNow := time.Now()
	timeBookingStart := timeNow.AddDate(0, 0, int(index))

	timeStr := timeBookingStart.Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)

	number, err := book_service.BookDoctor(id, t)
	if err != nil {
		logger.Error("Book doctor error", err)
		appG.Response(http.StatusOK, define.ErrorBooking, err.Error())
		return
	}
	appG.Response(http.StatusOK, define.SUCCESS, number)
}
