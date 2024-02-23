package services

import (
	"attendance/repository"
	"attendance/util"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ClassAttendanceService interface {
	GetClassAttendance(username string, data repository.GetClassAttendanceJSON) (int, repository.ErrorJSON, []repository.StudentAttendanceJSON)
}

type ClassAttendanceImpl struct {
	repository repository.Repository
}

func NewClassAttendanceImpl(repository repository.Repository) *ClassAttendanceImpl {
	return &ClassAttendanceImpl{repository: repository}
}

func (impl *ClassAttendanceImpl) GetClassAttendance(username string, data repository.GetClassAttendanceJSON) (int, repository.ErrorJSON, []repository.StudentAttendanceJSON) {
	user := impl.repository.GetUser(username)

	if user != nil && user.Role != "teacher" {
		zap.L().Info("Not authorized to get student attendance details")
		return http.StatusUnauthorized, repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1}, nil
	}

	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), data.Day, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), data.Day, 23, 59, 59)

	allStudentList := impl.repository.GetClassAttendance(data.Class, startDate, endDate)
	return http.StatusAccepted, repository.ErrorJSON{}, allStudentList
}
