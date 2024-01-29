package restHandler

import (
	"attendance/repository"
	auth "attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func AddNewUserHandler(w http.ResponseWriter, r *http.Request) {

	status, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	newUser := repository.User{
		UserID: uuid.New().String(),
		Email:  "",
	}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		zap.L().Error("Cannot decode json data for newUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	util.TrimSpacesFromStruct(&newUser)
	if newUser.IsNewUserDataMissing() {
		zap.L().Error("New user data is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = repository.AddNewUser(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zap.L().Info("New user added succesfully.")
	w.WriteHeader(http.StatusAccepted)
	return
}
