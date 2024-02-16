package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type GetStudentAttendanceService interface {
	GetStudentAttendance(username string, data repository.GetStudentAttendanceJSON, w http.ResponseWriter, r *http.Request)
}

type GetStudentAttendanceServiceImpl struct {
	repository repository.Repository
}

func NewGetStudentAttendanceServiceImpl(repository repository.Repository) *GetStudentAttendanceServiceImpl {
	return &GetStudentAttendanceServiceImpl{repository: repository}
}

func (impl *GetStudentAttendanceServiceImpl) GetStudentAttendance(username string, data repository.GetStudentAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := impl.repository.GetUser(username)

	if user == nil && user.Role != "student" {
		zap.L().Info("Not authorized to get student attendance details")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}
	allAttendances := impl.repository.GetStudentAttendance(username, data)
	json.NewEncoder(w).Encode(allAttendances)
}
