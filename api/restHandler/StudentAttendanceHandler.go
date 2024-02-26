package restHandler

import (
	auth "attendance/api/auth"
	"attendance/bean"
	"attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type StudentAttendanceHandler interface {
	GetStudentAttendance(w http.ResponseWriter, r *http.Request)
}

type StudentAttendanceImpl struct {
	studentAttendance services.StudentAttendanceService
	auth              auth.AuthToken
}

func NewStudentAttendanceImpl(studentAttendance services.StudentAttendanceService, auth auth.AuthToken) *StudentAttendanceImpl {
	return &StudentAttendanceImpl{studentAttendance: studentAttendance, auth: auth}
}

func ValidateStudentRequestData(data bean.GetStudentAttendanceJSON) (bool, string) {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func (impl *StudentAttendanceImpl) GetStudentAttendance(w http.ResponseWriter, r *http.Request) {
	status, username, _ := impl.auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newStudentAttendanceRequest := bean.GetStudentAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newStudentAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for student attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}
	flag, message := ValidateStudentRequestData(newStudentAttendanceRequest)
	if flag {
		zap.L().Info("Student attendance request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}

	flag, allAttendance := impl.studentAttendance.GetStudentAttendance(username, newStudentAttendanceRequest)
	if !flag {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(bean.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
	} else {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(allAttendance)
	}
}
