package repository

import (
	"attendance/util"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func GetUser(username string) *User {
	db := GetDB()
	user := &User{}

	err := db.Model(user).Where("username=?", username).Select()

	if err != nil {
		zap.L().Info(fmt.Sprintf("Username %v doesn't exist", username))
		return nil
	}
	return user
}

func AddNewUser(user *User) error {
	db := GetDB()
	_, err := db.Model(user).Insert()
	if err != nil {
		zap.L().Info("Error adding new user to DB.", zap.Error(err))
		return err
	}
	return nil
}

func GetCurrentStatus(username string) (bool, []Attendance) {
	t := util.GetCurrentIndianTime()
	startDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 0, 0, 0)
	endDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 23, 59, 59)

	db := GetDB()
	var attendances []Attendance
	db.Model(&attendances).
		Where("username = ?", username).
		Where("punch_in_date BETWEEN ? AND ?", startDate, endDate).
		Where("punch_out_date IS NULL").
		Select()
	if len(attendances) == 0 {
		return false, nil
	}
	return true, attendances
}

func AddNewPunchIn(username string) error {
	t := time.Now()
	_, currentTime := util.FormateDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	attendance := Attendance{
		AttendanceID: uuid.New().String(),
		PunchInDate:  currentTime,
		Username:     username,
	}
	db := GetDB()
	_, err := db.Model(&attendance).Insert()
	if err != nil {
		zap.L().Error("Cannot add new punch in of user "+username, zap.Error(err))
		return err
	}
	return nil
}

func AddNewPunchOut(username string, attendance Attendance) error {
	t := time.Now()
	_, currentTime := util.FormateDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	attendance.PunchOutDate = currentTime
	db := GetDB()
	_, err := db.Model(&attendance).Where("attendance_id = ?", attendance.AttendanceID).Column("punch_out_date").Update()
	if err != nil {
		zap.L().Error("Cannot add new punch out of user "+username, zap.Error(err))
		return err
	}
	return nil
}

func GetTeacherAttendance(username string) []Attendance {
	db := GetDB()
	var attendances []Attendance
	err := db.Model(&attendances).Where("username=?", username).Select()

	if err != nil {
		zap.L().Error("Cannot retrieve data for teacher "+username, zap.Error(err))
		return nil
	}

	return attendances
}
