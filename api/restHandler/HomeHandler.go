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

type HomeHandler interface {
	Home(w http.ResponseWriter, r *http.Request)
}

type HomeImpl struct {
	homeService services.HomeService
	auth        auth.AuthToken
}

func NewHomeImpl(homeService services.HomeService, auth auth.AuthToken) *HomeImpl {
	return &HomeImpl{homeService: homeService, auth: auth}
}

func ValidateRequestData(data bean.GetHomeJSON) (bool, string) {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func (impl *HomeImpl) Home(w http.ResponseWriter, r *http.Request) {
	status, username, role := impl.auth.VerifyToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newHomeRequest := bean.GetHomeJSON{}
	err := json.NewDecoder(r.Body).Decode(&newHomeRequest)

	if err != nil {
		zap.L().Error("Cannot decode json data for home request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}

	flag, message := ValidateRequestData(newHomeRequest)

	if flag {
		zap.L().Info("Home request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}
	if role == "student" {
		data := impl.homeService.StudentDashboardService(username, newHomeRequest)
		json.NewEncoder(w).Encode(data)
	} else if role == "teacher" {
		data := impl.homeService.TeacherDashboardService(username, newHomeRequest)
		json.NewEncoder(w).Encode(data)
	} else if role == "principal" {
		data, errorJSON := impl.homeService.PrincipalDashboardService(newHomeRequest)
		if data != nil {
			json.NewEncoder(w).Encode(data)
		} else {
			json.NewEncoder(w).Encode(errorJSON)
		}
	}
}
