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

type NewUserHandler interface {
	AddNewUser(w http.ResponseWriter, r *http.Request)
}

type NewUserImpl struct {
	newUserService services.NewUserService
	authService    auth.AuthService
}

func NewNewUserImpl(newUserService services.NewUserService, auth auth.AuthService) *NewUserImpl {
	return &NewUserImpl{newUserService: newUserService, authService: auth}
}

func (impl *NewUserImpl) AddNewUser(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)

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
