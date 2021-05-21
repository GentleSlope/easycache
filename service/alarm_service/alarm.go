package alarm_service

import (
	"easycache/models"
)

func InsertAlarmInfo(info models.AlarmInfo) (err error) {
	if err = models.InsertAlarmInfo(info); err != nil {
		return err
	}
	return nil
}

func GetAllAlarmInfos() (infos []models.AlarmInfo, err error) {
	infos, err = models.GetAllAlarmInfos()
	if err != nil {
		return infos, err
	}
	return infos, nil
}

func DeleteAlarmInfos() error {
	if err := models.DeleteAlarmInfos(); err != nil {
		return err
	}
	return nil
}
