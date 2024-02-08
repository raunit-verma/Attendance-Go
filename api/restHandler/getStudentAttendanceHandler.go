package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"attendance/util"
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

func GetStudentAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newStudentAttendanceRequest := repository.GetStudentAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newStudentAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for student attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}

	if ValidateStudentRequestData(newStudentAttendanceRequest) {
		zap.L().Info("Student attendance request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}

	services.GetStudentAttendanceService(username, newStudentAttendanceRequest, w, r)
	return
}
