package log_service

import (
	"easycache/models"
)

func InsertLogs(log models.LogInfo) (err error) {
	if err = models.InsertLogs(log); err != nil {
		return err
	}
	return nil
}

func GetAllLogs(ip string, port string, limit int) (logs []models.LogInfo, err error) {
	if logs, err = models.GetAllLogs(ip, port, limit); err != nil {
		return logs, err
	}
	return logs, nil
}

func DeleteLogs(ip string, port string) error {
	if err := models.DeleteLogs(ip, port); err != nil {
		return err
	}
	return nil
}
