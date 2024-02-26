package services

import (
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"time"

	"go.uber.org/zap"
)

type StudentAttendanceService interface {
	GetStudentAttendance(username string, data bean.GetStudentAttendanceJSON) (bool, []repository.Attendance)
}

type StudentAttendanceServiceImpl struct {
	repository repository.Repository
}

func NewStudentAttendanceServiceImpl(repository repository.Repository) *StudentAttendanceServiceImpl {
	return &StudentAttendanceServiceImpl{repository: repository}
}

func (impl *StudentAttendanceServiceImpl) GetStudentAttendance(username string, data bean.GetStudentAttendanceJSON) (bool, []repository.Attendance) {
	user := impl.repository.GetUser(username)

	if user != nil && user.Role != "student" {
		zap.L().Info("Not authorized to get student attendance details")
		return false, nil
	}

	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 1, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 31, 23, 59, 59)

	allAttendances := impl.repository.GetStudentAttendance(username, startDate, endDate)
	return true, allAttendances
}
