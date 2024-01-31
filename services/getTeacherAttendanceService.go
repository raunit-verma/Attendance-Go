package services

import (
	"attendance/repository"
	"encoding/json"
	"net/http"
)

func GetTeacherAttendanceService(username string, teacherId string, data repository.GetTeacherAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)

	if user.Role != "principal" && user.Role != "teacher" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user.Role == "teacher" {
		teacherId = user.Username
	}

	allAttendances := repository.GetTeacherAttendance(teacherId, data)
	json.NewEncoder(w).Encode(allAttendances)

}
