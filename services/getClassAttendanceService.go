package services

import (
	"attendance/repository"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func GetClassAttendanceService(username string, data repository.GetClassAttendanceJSON, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)

	if user.Role != "teacher" {
		zap.L().Info("Not authorized to get student attendance details")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	allStudentList := repository.GetClassAttendance(data)
	json.NewEncoder(w).Encode(allStudentList)
	return
}
