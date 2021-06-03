package book_service

import (
	"easycache/models"
	"time"
)

func SearchDoctors(department string, subject string, timeBookingStart time.Time) (dockers []models.Doctor, err error) {
	dockers, err = models.GetDoctorsBetween(department, subject, timeBookingStart)
	if err != nil {
		return dockers, err
	}
	return dockers, nil
}

func BookDoctor(id int64, start time.Time) (number int64, err error) {
	number, err = models.Booking(id, start)
	return
}
