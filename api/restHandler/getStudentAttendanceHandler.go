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

type GetStudentAttendanceHandler interface {
	GetStudentAttendance(w http.ResponseWriter, r *http.Request)
}

type GetStudentAttendanceImpl struct {
	getStudentAttendance services.GetStudentAttendanceService
}

func NewGetStudentAttendanceImpl(getStudentAttendance services.GetStudentAttendanceService) *GetStudentAttendanceImpl {
	return &GetStudentAttendanceImpl{getStudentAttendance: getStudentAttendance}
}

func ValidateStudentRequestData(data repository.GetStudentAttendanceJSON) (bool, string) {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func (impl *GetStudentAttendanceImpl) GetStudentAttendance(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
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
	flag, message := ValidateStudentRequestData(newStudentAttendanceRequest)
	if flag {
		zap.L().Info("Student attendance request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}

	impl.getStudentAttendance.GetStudentAttendance(username, newStudentAttendanceRequest, w, r)
	return
}
