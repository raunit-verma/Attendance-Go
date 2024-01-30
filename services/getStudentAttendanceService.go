package services

import (
	"attendance/repository"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func GetStudentAttendanceService(username string, data repository.GetStudentAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)

	if user == nil && user.Role != "student" {
		zap.L().Info("Not authorized to get student attendance details")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	allAttendances := repository.GetStudentAttendance(username, data)
	json.NewEncoder(w).Encode(allAttendances)
	return
}
