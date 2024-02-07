package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
)

func GetTeacherAttendanceService(username string, teacherId string, data repository.GetTeacherAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)

	if user.Role != "principal" && user.Role != "teacher" {
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}

	if user.Role == "teacher" {
		teacherId = user.Username
	}

	allAttendances := repository.GetTeacherAttendance(teacherId, data)
	json.NewEncoder(w).Encode(allAttendances)

}
