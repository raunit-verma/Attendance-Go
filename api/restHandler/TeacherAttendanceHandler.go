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

type TeacherAttendanceHandler interface {
	GetTeacherAttendance(w http.ResponseWriter, r *http.Request)
}

type TeacherAttendanceImpl struct {
	teacherAttendance services.TeacherAttendanceService
	auth              auth.AuthToken
}

func NewTeacherAttendanceImpl(teacherAttendance services.TeacherAttendanceService, auth auth.AuthToken) *TeacherAttendanceImpl {
	return &TeacherAttendanceImpl{teacherAttendance: teacherAttendance, auth: auth}
}

func ValidateTeacherRequestData(data bean.GetTeacherAttendanceJSON) (bool, string) {
	if data.ID == "" {
		zap.L().Info("Teacher id is null")
		return true, "Teacher id is null. "
	} else if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Month is not valid")
		return true, "Month is not valid. "
	} else if data.Year < 2020 || data.Year >= 2100 {
		zap.L().Info("Year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func (impl *TeacherAttendanceImpl) GetTeacherAttendance(w http.ResponseWriter, r *http.Request) {
	status, username, _ := impl.auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newTeacherAttendanceRequest := bean.GetTeacherAttendanceJSON{}
	err := json.NewDecoder(r.Body).Decode(&newTeacherAttendanceRequest)
	if err != nil {
		zap.L().Error("Cannot decode json data for teacher attendance request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}
	flag, message := ValidateTeacherRequestData(newTeacherAttendanceRequest)
	if flag {
		zap.L().Info("Teacher attendance request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}

	status, errorJSON, allAttendance := impl.teacherAttendance.GetTeacherAttendance(username, newTeacherAttendanceRequest.ID, newTeacherAttendanceRequest)
	w.WriteHeader(status)
	if status != http.StatusAccepted {
		json.NewEncoder(w).Encode(errorJSON)
	} else {
		json.NewEncoder(w).Encode(allAttendance)
	}
}
