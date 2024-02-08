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

func ValidateClassRequestData(data repository.GetClassAttendanceJSON) (bool, string) {
	if data.Class <= 0 || data.Class > 12 {
		zap.L().Info("Requested class is not valid")
		return true, "Class is not valid. "
	} else if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	} else if data.Day <= 0 || data.Day > 31 {
		zap.L().Info("Requested day is not valid")
		return true, "Day is not valid. "
	}
	return false, ""
}

func GetClassAttendanceHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		zap.L().Error("User not verified", zap.String("Code", "1"))
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newStudentAttendanceRequest := repository.GetClassAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newStudentAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for student attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}
	flag, message := ValidateClassRequestData(newStudentAttendanceRequest)
	if flag {
		zap.L().Info("Student attendance request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}

	services.GetClassAttendanceService(username, newStudentAttendanceRequest, w, r)
	return
}
