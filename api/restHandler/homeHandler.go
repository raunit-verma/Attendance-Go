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

func ValidateRequestData(data repository.GetHomeJSON) (bool, string) {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	status, username, role := auth.VerifyToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	newHomeRequest := repository.GetHomeJSON{}
	err := json.NewDecoder(r.Body).Decode(&newHomeRequest)

	if err != nil {
		zap.L().Error("Cannot decode json data for home request", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.CannotDecodePayload_Two, ErrorCode: 2})
		return
	}

	flag, message := ValidateRequestData(newHomeRequest)

	if flag {
		zap.L().Info("Home request data validation failed.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: message + util.RequestDataValidation_Five, ErrorCode: 5})
		return
	}
	if role == "student" {
		services.StudentDashboardService(w, r, username, newHomeRequest)
	} else if role == "teacher" {
		services.TeacherDashboardService(w, r, username, newHomeRequest)
	} else if role == "principal" {
		services.PrincipalDashboardService(w, r, newHomeRequest)
	}
}
