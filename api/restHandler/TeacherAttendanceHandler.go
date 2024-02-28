package restHandler

import (
	auth "attendance/api/auth"
	"attendance/bean"
	"attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"go.uber.org/zap"
)

type TeacherAttendanceHandler interface {
	GetTeacherAttendance(w http.ResponseWriter, r *http.Request)
}

type TeacherAttendanceImpl struct {
	teacherAttendance services.TeacherAttendanceService
	authService       auth.AuthService
}

func NewTeacherAttendanceImpl(teacherAttendance services.TeacherAttendanceService, auth auth.AuthService) *TeacherAttendanceImpl {
	return &TeacherAttendanceImpl{teacherAttendance: teacherAttendance, authService: auth}
}

func (impl *TeacherAttendanceImpl) GetTeacherAttendance(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)

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
