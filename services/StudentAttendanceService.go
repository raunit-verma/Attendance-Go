package services

import (
	"attendance/repository"

	"go.uber.org/zap"
)

type StudentAttendanceService interface {
	GetStudentAttendance(username string, data repository.GetStudentAttendanceJSON) (bool, []repository.Attendance)
}

type StudentAttendanceServiceImpl struct {
	repository repository.Repository
}

func NewStudentAttendanceServiceImpl(repository repository.Repository) *StudentAttendanceServiceImpl {
	return &StudentAttendanceServiceImpl{repository: repository}
}

func (impl *StudentAttendanceServiceImpl) GetStudentAttendance(username string, data repository.GetStudentAttendanceJSON) (bool, []repository.Attendance) {
	user := impl.repository.GetUser(username)

	if user != nil && user.Role != "student" {
		zap.L().Info("Not authorized to get student attendance details")
		return false, nil
	}

	allAttendances := impl.repository.GetStudentAttendance(username, data)
	return true, allAttendances
}
