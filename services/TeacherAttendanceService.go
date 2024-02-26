package services

import (
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"net/http"
	"time"
)

type TeacherAttendanceService interface {
	GetTeacherAttendance(username string, teacherId string, data bean.GetTeacherAttendanceJSON) (int, bean.ErrorJSON, []repository.Attendance)
}

type TeacherAttendanceServiceImpl struct {
	repository repository.Repository
}

func NewTeacherAttendanceServiceImpl(repository repository.Repository) *TeacherAttendanceServiceImpl {
	return &TeacherAttendanceServiceImpl{repository: repository}
}

func (impl *TeacherAttendanceServiceImpl) GetTeacherAttendance(username string, teacherId string, data bean.GetTeacherAttendanceJSON) (int, bean.ErrorJSON, []repository.Attendance) {
	user := impl.repository.GetUser(username)

	if user == nil || (user.Role != "principal" && user.Role != "teacher") {
		return http.StatusUnauthorized, bean.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One}, nil
	}

	if user.Role == "teacher" {
		teacherId = user.Username
	}

	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 1, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 31, 23, 59, 59)

	allAttendances := impl.repository.GetTeacherAttendance(teacherId, startDate, endDate)
	return http.StatusAccepted, bean.ErrorJSON{}, allAttendances
}
