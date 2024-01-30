package services

import (
	"attendance/repository"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func GetStudentsAttendanceService(username string, data repository.GetStudentAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)

	if user.Role != "teacher" {
		zap.L().Info("Not authorized to get student attendance details")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	allStudentList := repository.GetStudentAttendances(data)
	json.NewEncoder(w).Encode(allStudentList)
	return
}
