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
		return
	}

	err := repository.AddNewUser(&newUser)
	pgErr, ok := err.(pg.Error)
	if err != nil {
		if ok && pgErr.Field('C') == "23505" {
			json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 4, Message: util.Four})
			zap.L().Error("Username already exist", zap.String("username", newUser.Username))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zap.L().Info("New user added succesfully.")
	w.WriteHeader(http.StatusAccepted)
	return
}
