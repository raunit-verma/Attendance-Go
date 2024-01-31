package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func ValidateStudentRequestData(data repository.GetStudentAttendanceJSON) bool {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true
	}
	return false
}

func GetStudentsAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	newStudentsAttendanceRequest := repository.GetStudentAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newStudentsAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for student attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ValidateStudentRequestData(newStudentsAttendanceRequest) {
		zap.L().Info("Student attendance request data validation failed.")
		return
	}

	services.GetStudentAttendanceService(username, newStudentsAttendanceRequest, w, r)
	return
}