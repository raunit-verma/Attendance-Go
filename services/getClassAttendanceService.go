package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type GetClassAttendanceService interface {
	GetClassAttendance(username string, data repository.GetClassAttendanceJSON, w http.ResponseWriter, r *http.Request)
}

type GetClassAttendanceImpl struct {
	repository repository.Repository
}

func NewGetClassAttendanceImpl(repository repository.Repository) *GetClassAttendanceImpl {
	return &GetClassAttendanceImpl{repository: repository}
}

func (impl *GetClassAttendanceImpl) GetClassAttendance(username string, data repository.GetClassAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := impl.repository.GetUser(username)

	if user.Role != "teacher" {
		zap.L().Info("Not authorized to get student attendance details")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	allStudentList := impl.repository.GetClassAttendance(data)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(allStudentList)
	return
}
