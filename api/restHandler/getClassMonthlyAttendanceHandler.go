package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func ValidateClassRequestData(data repository.GetClassAttendanceJSON) bool {
	if data.Class <= 0 || data.Class > 12 {
		zap.L().Info("Requested class is not valid")
		return true
	} else if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true
	} else if data.Day <= 0 || data.Day > 31 {
		zap.L().Info("Requested day is not valid")
		return true
	}
	return false
}

func GetClassAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	newStudentAttendanceRequest := repository.GetClassAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newStudentAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for student attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ValidateClassRequestData(newStudentAttendanceRequest) {
		zap.L().Info("Student attendance request data validation failed.")
		return
	}

	services.GetClassAttendanceService(username, newStudentAttendanceRequest, w, r)
	return
}
