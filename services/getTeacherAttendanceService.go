package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
)

type GetTeacherAttendanceService interface {
	GetTeacherAttendance(username string, teacherId string, data repository.GetTeacherAttendanceJSON, w http.ResponseWriter, r *http.Request)
}

type GetTeacherAttendanceServiceImpl struct {
	repository repository.Repository
}

func NewGetTeacherAttendanceServiceImpl(repository repository.Repository) *GetTeacherAttendanceServiceImpl {
	return &GetTeacherAttendanceServiceImpl{repository: repository}
}

func (impl *GetTeacherAttendanceServiceImpl) GetTeacherAttendance(username string, teacherId string, data repository.GetTeacherAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := impl.repository.GetUser(username)

	if user.Role != "principal" && user.Role != "teacher" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}

	if user.Role == "teacher" {
		teacherId = user.Username
	}

	allAttendances := impl.repository.GetTeacherAttendance(teacherId, data)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(allAttendances)
}
