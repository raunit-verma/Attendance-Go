package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func AddNewUserHandler(w http.ResponseWriter, r *http.Request) {

	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	newUser := repository.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		zap.L().Error("Cannot decode json data for newUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	services.AddNewUserService(newUser, username, w, r)

	return
}
