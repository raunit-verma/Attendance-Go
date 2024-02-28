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
	auth        auth.AuthService
}

func NewHomeImpl(homeService services.HomeService, auth auth.AuthService) *HomeImpl {
	return &HomeImpl{homeService: homeService, auth: auth}
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

	flag, message := ValidateHomeRequestData(newHomeRequest)

	if flag {
		zap.L().Info("Home request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}
	if role == util.STUDENT {
		data := impl.homeService.StudentDashboard(username, newHomeRequest)
		json.NewEncoder(w).Encode(data)
	} else if role == util.TEACHER {
		data := impl.homeService.TeacherDashboard(username, newHomeRequest)
		json.NewEncoder(w).Encode(data)
	} else if role == util.PRINCIPAL {
		data, errorJSON := impl.homeService.PrincipalDashboard(newHomeRequest)
		if data != nil {
			json.NewEncoder(w).Encode(data)
		} else {
			json.NewEncoder(w).Encode(errorJSON)
		}
	}
}
