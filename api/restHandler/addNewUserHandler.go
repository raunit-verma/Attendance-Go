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

type AddNewUserHandler interface {
	AddNewUser(w http.ResponseWriter, r *http.Request)
}

type AddNewUserImpl struct {
	addNewUserService services.AddNewUserService
}

func NewAddNewUserImpl(addNewUserService services.AddNewUserService) *AddNewUserImpl {
	return &AddNewUserImpl{addNewUserService: addNewUserService}
}

func (impl *AddNewUserImpl) AddNewUser(w http.ResponseWriter, r *http.Request) {

	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}

	newUser := repository.User{}

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		zap.L().Error("Cannot decode json data for newUser", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 2, Message: util.CannotDecodePayload_Two})
		return
	}

	impl.addNewUserService.AddNewUser(newUser, username, w, r)

	return
}
