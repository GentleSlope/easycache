package monitor_service

import "easycache/models"

func InsertInfo(info models.MonitorInfo) (err error) {
	if err := models.InsertInfo(info); err != nil {
		return err
	}
	return nil
}

func GetMonitorInfo(ip string, port int64, limit int) (infos []models.MonitorInfo, err error) {
	if infos, err = models.GetMonitorInfo(ip, port, limit); err != nil {
		return infos, err
	}
	return infos, nil
}
