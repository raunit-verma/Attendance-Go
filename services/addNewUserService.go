package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

func AddNewUserService(newUser repository.User, username string, w http.ResponseWriter, r *http.Request) {
	util.TrimSpacesFromStruct(&newUser)
	if newUser.IsNewUserDataMissing(w, r) {
		zap.L().Error("New user data is missing")
		return
	}

	user := repository.GetUser(username)
	if user.Role != "principal" {
		zap.L().Warn("Unauthorized to add new user")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	err := repository.AddNewUser(&newUser)
	pgErr, ok := err.(pg.Error)
	if err != nil {
		if ok && pgErr.Field('C') == "23505" {
			zap.L().Error("Username already exist", zap.String("username", newUser.Username))
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 4, Message: util.UsernameTaken_Four})
			return
		}
		w.WriteHeader(util.InternalServererror)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7})
		return
	}

	zap.L().Info("New user added succesfully.")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.Success_Eight + " Added new user " + username, ErrorCode: 8})
	return
}
