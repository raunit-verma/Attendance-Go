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

type ClassAttendanceHandler interface {
	GetClassAttendance(w http.ResponseWriter, r *http.Request)
}

type ClassAttendanceImpl struct {
	classAttendance services.ClassAttendanceService
	authService     auth.AuthService
}

func NewClassAttendanceImpl(classAttendance services.ClassAttendanceService, auth auth.AuthService) *ClassAttendanceImpl {
	return &ClassAttendanceImpl{classAttendance: classAttendance, authService: auth}
}

func (impl *ClassAttendanceImpl) GetClassAttendance(w http.ResponseWriter, r *http.Request) {
	status, username, _ := impl.authService.VerifyToken(r)
	if status != http.StatusAccepted {
		zap.L().Error("User not verified", zap.String("Code", "1"))
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newStudentAttendanceRequest := bean.GetClassAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newStudentAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for student attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}
	flag, message := ValidateClassRequestData(newStudentAttendanceRequest)
	if flag {
		zap.L().Info("Student attendance request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}

	status, errorJSON, allStudentList := impl.classAttendance.GetClassAttendance(username, newStudentAttendanceRequest)
	w.WriteHeader(status)
	if status == http.StatusAccepted {
		json.NewEncoder(w).Encode(allStudentList)
		return
	}
	json.NewEncoder(w).Encode(errorJSON)

}
