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

type NewUserHandler interface {
	AddNewUser(w http.ResponseWriter, r *http.Request)
}

type NewUserImpl struct {
	newUserService services.NewUserService
	auth           auth.AuthToken
}

func NewNewUserImpl(newUserService services.NewUserService, auth auth.AuthToken) *NewUserImpl {
	return &NewUserImpl{newUserService: newUserService, auth: auth}
}

func (impl *NewUserImpl) AddNewUser(w http.ResponseWriter, r *http.Request) {
	status, username, _ := impl.auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}

	newUser := bean.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		zap.L().Error("Cannot decode json data for newUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(bean.ErrorJSON{ErrorCode: 2, Message: util.CannotDecodePayload_Two})
		return
	}
	status, errorJSON := impl.newUserService.AddNewUser(newUser, username)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorJSON)
}
