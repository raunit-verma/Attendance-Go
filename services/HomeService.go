package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
	"time"
)

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

func TeacherDashboardService(w http.ResponseWriter, r *http.Request, username string, requestData repository.GetHomeJSON) {
	data := repository.GetTeacherAttendanceJSON{ID: username, Month: requestData.Month, Year: requestData.Year}
	allAttendance := repository.GetTeacherAttendance(username, data)
	monthlyAttendance, duration := getMonthlyAttendance(allAttendance)

	json.NewEncoder(w).Encode(repository.DashboardJSON{MonthlyAttendance: monthlyAttendance[:], Hour: int(duration.Hours()), Minute: int(duration.Minutes()) % 60, Second: int(duration.Seconds()) % 60})
}

func StudentDashboardService(w http.ResponseWriter, r *http.Request, username string, requestData repository.GetHomeJSON) {
	data := repository.GetStudentAttendanceJSON{Month: requestData.Month, Year: requestData.Year}
	allAttendance := repository.GetStudentAttendance(username, data)
	monthlyAttendance, duration := getMonthlyAttendance(allAttendance)

	json.NewEncoder(w).Encode(repository.DashboardJSON{MonthlyAttendance: monthlyAttendance[:], Hour: int(duration.Hours()), Minute: int(duration.Minutes()) % 60, Second: int(duration.Seconds()) % 60})
}

func PrincipalDashboardService(w http.ResponseWriter, r *http.Request, requestData repository.GetHomeJSON) {
	totalStudentPresent, totalTeacherPresent, totalStudent, totalTeacher := repository.GetDailyStats(requestData)
	if totalStudentPresent == -1 {
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 7, Message: util.DBError_Seven})
	}
	json.NewEncoder(w).Encode(map[string]int{"totalStudentPresent": totalStudentPresent, "totalTeacherPresent": totalTeacherPresent, "totalStudent": totalStudent, "totalTeacher": totalTeacher})
}
