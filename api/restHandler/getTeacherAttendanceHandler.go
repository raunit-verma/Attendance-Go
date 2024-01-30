package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func ValidateTeacherRequestData(data repository.GetTeacherAttendanceJSON) bool {
	if data.ID == "" {
		zap.L().Info("Teacher id is null")
		return true
	} else if data.Month < 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true
	}
	return false
}

func GetTeacherAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	newTeacherAttendanceRequest := repository.GetTeacherAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newTeacherAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for teacher attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ValidateTeacherRequestData(newTeacherAttendanceRequest) {
		zap.L().Info("Teacher attendance request data validation failed.")
		return
	}

	services.GetTeacherAttendanceService(username, newTeacherAttendanceRequest.ID, w, r)
	return
}
