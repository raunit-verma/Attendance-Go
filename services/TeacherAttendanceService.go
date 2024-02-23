package services

import (
	"attendance/repository"
	"attendance/util"
	"net/http"
)

type TeacherAttendanceService interface {
	GetTeacherAttendance(username string, teacherId string, data repository.GetTeacherAttendanceJSON) (int, repository.ErrorJSON, []repository.Attendance)
}

type TeacherAttendanceServiceImpl struct {
	repository repository.Repository
}

func NewTeacherAttendanceServiceImpl(repository repository.Repository) *TeacherAttendanceServiceImpl {
	return &TeacherAttendanceServiceImpl{repository: repository}
}

func (impl *TeacherAttendanceServiceImpl) GetTeacherAttendance(username string, teacherId string, data repository.GetTeacherAttendanceJSON) (int, repository.ErrorJSON, []repository.Attendance) {
	user := impl.repository.GetUser(username)

	if user == nil || (user.Role != "principal" && user.Role != "teacher") {
		return http.StatusUnauthorized, repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One}, nil
	}

	if user.Role == "teacher" {
		teacherId = user.Username
	}

	allAttendances := impl.repository.GetTeacherAttendance(teacherId, data)
	return http.StatusAccepted, repository.ErrorJSON{}, allAttendances
}
