package models

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"time"
)

type Doctor struct {
	DoctorId       int64  `gorm:"primary_key" json:"doctor_id"`
	DoctorName     string `json:"doctor_name"`
	JobTitle       string `json:"job_title"`
	JobTitleId     int64  `json:"job_title_id"`
	DepartmentName string `json:"department_name"`
	DepartmentId   int64  `json:"department_id"`
	SubjectName    string `json:"subject_name"`
	SubjectId      int64  `json:"subject_id"`
	Sex            string `json:"sex"`
}

type DoctorBooking struct {
	Id       int64     `gorm:"primary_key" json:"id"`
	DoctorId int64     `json:"doctor_id"`
	Date     time.Time `json:"date"`
	Limit    int64     `json:"limit"`
	Number   int64     `json:"number"`
}

// ...
func GetDoctors() (doctors []Doctor, err error) {
	err = db.Where("1 = 1").Find(&doctors).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return doctors, err
	}
	return doctors, nil
}

// 为所有医生都插入一个指定时间的预定表
func InsertBookings(start time.Time) (err error) {
	doctors, err := GetDoctors()
	if err != nil {
		return err
	}
	for _, v := range doctors {
		booking := DoctorBooking{
			DoctorId: v.DoctorId,
			Date:     start,
			Limit:    5,
		}
		if err := db.Create(&booking).Error; err != nil {
			return err
		}
	}
	return nil
}

// GetDoctorsBetween ...
func GetDoctorsBetween(department string, subject string, start time.Time) (doctors []Doctor, err error) {
	var bookings []DoctorBooking
	err = db.Where("date = ?", start.String()).Find(&bookings).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 没有的话就直接为所有医生插入一个
	if len(bookings) == 0 {
		err = InsertBookings(start)
		if err != nil {
			return doctors, err
		}
		err = db.Where("date = ?", start.String()).Find(&bookings).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	var doctorIds = make([]int64, len(bookings))
	for _, v := range bookings {
		// 还未预约满的医生
		if v.Number < v.Limit {
			doctorIds = append(doctorIds, v.DoctorId)
		}
	}
	err = db.Where("department_name = ? AND subject_name = ? AND doctor_id IN ( ? )", department, subject, doctorIds).Find(&doctors).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return doctors, nil
}

func Booking(id int64, date time.Time) (number int64, err error) {
	var doctor Doctor
	err = db.Where("doctor_id = ?", id).First(&doctor).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return number, err
	}
	if doctor.DoctorId == 0 {
		return number, errors.New("没有id对应的医生！")
	}

	var booking DoctorBooking
	err = db.Where("doctor_id = ? AND date = ?", id, date.String()).First(&booking).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return number, err
	}
	if doctor.DoctorId == 0 {
		return number, errors.New("没有对应的订单！")
	}
	booking.Number = booking.Number + 1
	if err = db.Save(&booking).Error; err != nil {
		return number, err
	}

	return booking.Number, nil
}
