package services

import (
	"attendance/repository"
	"attendance/util"
	"time"
)

type HomeService interface {
	TeacherDashboardService(username string, requestData repository.GetHomeJSON) repository.DashboardJSON
	StudentDashboardService(username string, requestData repository.GetHomeJSON) repository.DashboardJSON
	PrincipalDashboardService(requestData repository.GetHomeJSON) (map[string]int, repository.ErrorJSON)
}

type HomeServiceImpl struct {
	repository repository.Repository
}

func NewHomeServiceImpl(repository repository.Repository) *HomeServiceImpl {
	return &HomeServiceImpl{repository: repository}
}

func getMonthlyAttendance(allAttendance []repository.Attendance) ([32]bool, time.Duration) {
	var monthlyAttendance [32]bool
	var totalSeconds int64 = 0
	for i := 0; i < len(allAttendance); i++ {
		if allAttendance[i].PunchOutDate.After(allAttendance[i].PunchInDate) {
			monthlyAttendance[allAttendance[i].PunchInDate.Day()] = true
			totalSeconds += allAttendance[i].PunchOutDate.Unix() - allAttendance[i].PunchInDate.Unix()
		}
	}
	duration := time.Duration(float64(totalSeconds) * 1e9)
	return monthlyAttendance, duration
}

func (impl *HomeServiceImpl) TeacherDashboardService(username string, requestData repository.GetHomeJSON) repository.DashboardJSON {
	data := repository.GetTeacherAttendanceJSON{ID: username, Month: requestData.Month, Year: requestData.Year}
	allAttendance := impl.repository.GetTeacherAttendance(username, data)
	monthlyAttendance, duration := getMonthlyAttendance(allAttendance)
	return repository.DashboardJSON{MonthlyAttendance: monthlyAttendance[:], Hour: int(duration.Hours()), Minute: int(duration.Minutes()) % 60, Second: int(duration.Seconds()) % 60}
}

func (impl *HomeServiceImpl) StudentDashboardService(username string, requestData repository.GetHomeJSON) repository.DashboardJSON {
	data := repository.GetStudentAttendanceJSON{Month: requestData.Month, Year: requestData.Year}
	allAttendance := impl.repository.GetStudentAttendance(username, data)
	monthlyAttendance, duration := getMonthlyAttendance(allAttendance)
	return repository.DashboardJSON{MonthlyAttendance: monthlyAttendance[:], Hour: int(duration.Hours()), Minute: int(duration.Minutes()) % 60, Second: int(duration.Seconds()) % 60}
}

func (impl *HomeServiceImpl) PrincipalDashboardService(requestData repository.GetHomeJSON) (map[string]int, repository.ErrorJSON) {
	totalStudentPresent, totalTeacherPresent, totalStudent, totalTeacher := impl.repository.GetDailyStats(requestData)
	if totalStudentPresent == -1 {
		return nil, repository.ErrorJSON{ErrorCode: 7, Message: util.DBError_Seven}
	}
	return map[string]int{"totalStudentPresent": totalStudentPresent, "totalTeacherPresent": totalTeacherPresent, "totalStudent": totalStudent, "totalTeacher": totalTeacher}, repository.ErrorJSON{}
}
